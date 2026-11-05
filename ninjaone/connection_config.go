package ninjaone

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

const clientCacheKey = "ninjaone.client"

type ninjaoneConfig struct {
	// ClientID is the client ID of an API Services (machine-to-machine) client app.
	ClientID *string `hcl:"client_id"`

	// ClientSecret is that client app's secret.
	ClientSecret *string `hcl:"client_secret"`

	// BaseURL is the regional instance, such as https://oc.ninjarmm.com.
	BaseURL *string `hcl:"base_url"`

	// OrganizationID limits every table to one NinjaOne organization. It is a
	// string so the value survives connection generators that only write strings.
	OrganizationID *string `hcl:"organization_id"`
}

// ConfigInstance returns a new, empty connection configuration struct.
func ConfigInstance() interface{} {
	return &ninjaoneConfig{}
}

// GetConfig retrieves and casts the connection configuration.
func GetConfig(connection *plugin.Connection) ninjaoneConfig {
	if connection == nil {
		return ninjaoneConfig{}
	}
	config, _ := connection.GetConfig().(ninjaoneConfig)
	return config
}

func setting(value *string, envName string) string {
	if value != nil && *value != "" {
		return *value
	}
	return os.Getenv(envName)
}

func parseOrganizationID(value string) (*int64, error) {
	if value == "" {
		return nil, nil
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return nil, fmt.Errorf("organization_id must be a positive whole number, got %q", value)
	}
	return &id, nil
}

// connect returns the connection's client, building it once so its token is reused.
func connect(ctx context.Context, d *plugin.QueryData) (*client, error) {
	if cached, ok := d.ConnectionCache.Get(ctx, clientCacheKey); ok {
		return cached.(*client), nil
	}
	config := GetConfig(d.Connection)
	clientID := setting(config.ClientID, "NINJAONE_CLIENT_ID")
	clientSecret := setting(config.ClientSecret, "NINJAONE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("client_id and client_secret must be set in the connection config, or NINJAONE_CLIENT_ID and NINJAONE_CLIENT_SECRET in the environment")
	}
	baseURL := setting(config.BaseURL, "NINJAONE_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	organizationID, err := parseOrganizationID(setting(config.OrganizationID, "NINJAONE_ORGANIZATION_ID"))
	if err != nil {
		return nil, err
	}
	c := newClient(clientID, clientSecret, baseURL, organizationID)
	if err := d.ConnectionCache.Set(ctx, clientCacheKey, c); err != nil {
		plugin.Logger(ctx).Warn("ninjaone.connect", "cache_error", err)
	}
	return c, nil
}
