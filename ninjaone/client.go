package ninjaone

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	defaultBaseURL = "https://app.ninjarmm.com"
	maxPageSize    = 1000
	maxAttempts    = 5
	requestTimeout = 60 * time.Second
	maxBackoff     = 30 * time.Second
	errorBodyLimit = 512
)

// apiError is a non-2xx response from the NinjaOne API.
type apiError struct {
	Path   string
	Status int
	Body   string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("GET %s returned %d: %s", e.Path, e.Status, e.Body)
}

func isNotFound(err error) bool {
	var apiErr *apiError
	return errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound
}

// client reads the NinjaOne public API. When organizationID is set, every
// read it serves is limited to that one organization.
type client struct {
	baseURL        string
	http           *http.Client
	organizationID *int64
	wait           func(context.Context, time.Duration) error
}

func newClient(clientID, clientSecret, baseURL string, organizationID *int64) *client {
	baseURL = strings.TrimRight(baseURL, "/")
	credentials := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     baseURL + "/ws/oauth/token",
		Scopes:       []string{"monitoring"},
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	// The client outlives any one query, so a query's context must not own token refresh.
	tokenContext := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Timeout: requestTimeout})
	httpClient := credentials.Client(tokenContext)
	httpClient.Timeout = requestTimeout
	return &client{baseURL: baseURL, http: httpClient, organizationID: organizationID, wait: sleep}
}

func sleep(ctx context.Context, delay time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		return nil
	}
}

func retryable(status int) bool {
	switch status {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

func retryDelay(retryAfter string, attempt int) time.Duration {
	if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds >= 0 {
		return min(time.Duration(seconds)*time.Second, maxBackoff)
	}
	return min(time.Second<<(attempt-1), maxBackoff)
}

// get decodes one API response into out, retrying throttling and gateway errors.
// Retries live here, not in the SDK's hydrate retry, because a hydrate retry
// would restart a list and stream its rows twice.
func (c *client) get(ctx context.Context, path string, query url.Values, out any) error {
	target := c.baseURL + path
	if len(query) > 0 {
		target += "?" + query.Encode()
	}
	for attempt := 1; ; attempt++ {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
		if err != nil {
			return fmt.Errorf("build request for %s: %w", path, err)
		}
		request.Header.Set("Accept", "application/json")
		response, err := c.http.Do(request)
		if err != nil {
			return fmt.Errorf("GET %s: %w", path, err)
		}
		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if response.StatusCode >= 200 && response.StatusCode < 300 {
			if err := json.Unmarshal(body, out); err != nil {
				return fmt.Errorf("decode %s: %w", path, err)
			}
			return nil
		}
		if !retryable(response.StatusCode) || attempt == maxAttempts {
			return &apiError{Path: path, Status: response.StatusCode, Body: truncate(string(body))}
		}
		if err := c.wait(ctx, retryDelay(response.Header.Get("Retry-After"), attempt)); err != nil {
			return err
		}
	}
}

func truncate(body string) string {
	if len(body) > errorBodyLimit {
		return body[:errorBodyLimit]
	}
	return body
}

// deviceFilter builds the df parameter: the connection's organization, narrowed
// to one device when asked. Empty means no filter.
func (c *client) deviceFilter(deviceID *int64) string {
	var terms []string
	if c.organizationID != nil {
		terms = append(terms, fmt.Sprintf("org=%d", *c.organizationID))
	}
	if deviceID != nil {
		terms = append(terms, fmt.Sprintf("id=%d", *deviceID))
	}
	return strings.Join(terms, " AND ")
}

func (c *client) inScope(organizationID int64) bool {
	return c.organizationID == nil || *c.organizationID == organizationID
}

func pageSize(limit int64) int {
	if limit > 0 && limit < maxPageSize {
		return int(limit)
	}
	return maxPageSize
}

func filterQuery(filter string) url.Values {
	query := url.Values{}
	if filter != "" {
		query.Set("df", filter)
	}
	return query
}

// listCursor walks a /v2/queries report. each returns false to stop early.
func (c *client) listCursor(ctx context.Context, path, filter string, limit int64, each func(json.RawMessage) (bool, error)) error {
	query := filterQuery(filter)
	query.Set("pageSize", strconv.Itoa(pageSize(limit)))
	for {
		var page struct {
			Cursor *struct {
				Name string `json:"name"`
			} `json:"cursor"`
			Results []json.RawMessage `json:"results"`
		}
		if err := c.get(ctx, path, query, &page); err != nil {
			return err
		}
		for _, row := range page.Results {
			more, err := each(row)
			if err != nil || !more {
				return err
			}
		}
		if len(page.Results) == 0 || page.Cursor == nil || page.Cursor.Name == "" {
			return nil
		}
		query.Set("cursor", page.Cursor.Name)
	}
}

// listAfter walks an endpoint that pages by the last id seen.
func (c *client) listAfter(ctx context.Context, path, filter string, limit int64, each func(json.RawMessage) (bool, error)) error {
	query := filterQuery(filter)
	query.Set("pageSize", strconv.Itoa(pageSize(limit)))
	var last int64
	for {
		var rows []json.RawMessage
		if err := c.get(ctx, path, query, &rows); err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		previous := last
		for _, row := range rows {
			var key struct {
				ID int64 `json:"id"`
			}
			if err := json.Unmarshal(row, &key); err != nil {
				return fmt.Errorf("decode %s id: %w", path, err)
			}
			last = max(last, key.ID)
			more, err := each(row)
			if err != nil || !more {
				return err
			}
		}
		if last == previous {
			return fmt.Errorf("%s did not advance past id %d", path, last)
		}
		query.Set("after", strconv.FormatInt(last, 10))
	}
}

// requireOrganization fails when the connection's organization does not exist,
// so a mistyped or removed id can never read as a tenant with no devices.
func (c *client) requireOrganization(ctx context.Context) error {
	if c.organizationID == nil {
		return nil
	}
	var organization struct {
		ID int64 `json:"id"`
	}
	return c.get(ctx, fmt.Sprintf("/v2/organization/%d", *c.organizationID), nil, &organization)
}

// scopedDevices lists the organization's device ids, or nil on an unscoped connection.
func (c *client) scopedDevices(ctx context.Context) (map[int64]bool, error) {
	if c.organizationID == nil {
		return nil, nil
	}
	if err := c.requireOrganization(ctx); err != nil {
		return nil, err
	}
	devices := map[int64]bool{}
	err := c.listAfter(ctx, "/v2/devices-detailed", c.deviceFilter(nil), -1, func(raw json.RawMessage) (bool, error) {
		var device struct {
			ID             int64 `json:"id"`
			OrganizationID int64 `json:"organizationId"`
		}
		if err := json.Unmarshal(raw, &device); err != nil {
			return false, fmt.Errorf("decode device: %w", err)
		}
		if c.inScope(device.OrganizationID) {
			devices[device.ID] = true
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return devices, nil
}

// onlyDevices drops rows that are not on an allowed device, as a second check
// behind the df filter. A nil allow-list means an unscoped connection.
func onlyDevices(allowed map[int64]bool, next func(json.RawMessage) (bool, error)) func(json.RawMessage) (bool, error) {
	if allowed == nil {
		return next
	}
	return func(raw json.RawMessage) (bool, error) {
		var row struct {
			DeviceID *int64 `json:"deviceId"`
		}
		if err := json.Unmarshal(raw, &row); err != nil {
			return false, fmt.Errorf("decode row device: %w", err)
		}
		if row.DeviceID == nil || !allowed[*row.DeviceID] {
			return true, nil
		}
		return next(raw)
	}
}

// listActivities walks /v2/activities, newest first, by the oldest id seen.
func (c *client) listActivities(ctx context.Context, filter string, newerThan *int64, limit int64, each func(json.RawMessage) (bool, error)) error {
	query := filterQuery(filter)
	query.Set("pageSize", strconv.Itoa(pageSize(limit)))
	if newerThan != nil {
		query.Set("newerThan", strconv.FormatInt(*newerThan, 10))
	}
	for {
		var page struct {
			Activities []json.RawMessage `json:"activities"`
		}
		if err := c.get(ctx, "/v2/activities", query, &page); err != nil {
			return err
		}
		if len(page.Activities) == 0 {
			return nil
		}
		var oldest int64
		for _, row := range page.Activities {
			var key struct {
				ID int64 `json:"id"`
			}
			if err := json.Unmarshal(row, &key); err != nil {
				return fmt.Errorf("decode activity id: %w", err)
			}
			if oldest == 0 || key.ID < oldest {
				oldest = key.ID
			}
			more, err := each(row)
			if err != nil || !more {
				return err
			}
		}
		query.Set("olderThan", strconv.FormatInt(oldest, 10))
	}
}
