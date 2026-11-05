STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo

install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/local/ninjaone/ninjaone.plugin -tags "$(BUILD_TAGS)" *.go

# Build the plugin as a SQLite loadable extension (steampipe_sqlite_ninjaone.so)
# using turbot/steampipe-sqlite with a replace directive pointing at this repo.
sqlite:
	rm -rf build/steampipe-sqlite && mkdir -p build
	git clone --quiet --depth 1 https://github.com/turbot/steampipe-sqlite.git build/steampipe-sqlite
	cd build/steampipe-sqlite && rm -rf work && mkdir -p work && rsync -a --exclude='.git' . work/
	cd build/steampipe-sqlite/work && \
		go run generate/generator.go templates . ninjaone "" github.com/konfirmity/steampipe-plugin-ninjaone && \
		echo 'replace github.com/konfirmity/steampipe-plugin-ninjaone => $(CURDIR)' >> go.mod && \
		go mod tidy && \
		$(MAKE) -f out/Makefile build
	cp build/steampipe-sqlite/work/steampipe_sqlite_ninjaone.so .
	@echo "Built ./steampipe_sqlite_ninjaone.so"

# Build the plugin as a standalone Postgres FDW extension
# (steampipe_postgres_ninjaone.so + .control + .sql in build/postgres/) using
# turbot/steampipe-postgres-fdw. Requires pg_config on PATH.
# Pinned deliberately. The FDW's main branch resolves steampipe-plugin-sdk transitively through
# steampipe/v2, so an unpinned clone can change SDK major version between builds with no change
# here; when it does, the build fails with a PluginFunc type mismatch against this plugin's SDK.
# v2.3.0-rc.0 is the tag whose resolved SDK is v6, matching go.mod.
FDW_VERSION ?= v2.3.0-rc.0

postgres:
	rm -rf build/steampipe-postgres-fdw && mkdir -p build
	git clone --quiet --depth 1 --branch $(FDW_VERSION) https://github.com/turbot/steampipe-postgres-fdw.git build/steampipe-postgres-fdw
	cd build/steampipe-postgres-fdw && \
		$(MAKE) prebuild.go && \
		rm -rf work && mkdir -p work && rsync -a --exclude='.git' . work/
	cd build/steampipe-postgres-fdw/work && \
		go run generate/generator.go templates . ninjaone "" github.com/konfirmity/steampipe-plugin-ninjaone && \
		echo 'replace github.com/konfirmity/steampipe-plugin-ninjaone => $(CURDIR)' >> go.mod && \
		go mod tidy && \
		$(MAKE) -C ./fdw clean plugin=ninjaone && \
		$(MAKE) -C ./fdw go plugin=ninjaone && \
		$(MAKE) -C ./fdw plugin=ninjaone && \
		$(MAKE) -C ./fdw standalone plugin=ninjaone
	rm -rf build/postgres && mkdir -p build/postgres
	cp -a build/steampipe-postgres-fdw/work/build-$(shell uname)/* build/postgres/
	@echo "Built build/postgres/ (steampipe_postgres_ninjaone extension)"

clean:
	rm -rf build steampipe_sqlite_ninjaone.so
