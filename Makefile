.PHONY: build package run stop run-client run-server stop-client stop-server restart restart-server restart-client start-docker clean-dist clean nuke check-style check-client-style check-server-style check-unit-tests test dist setup-mac prepare-enteprise run-client-tests setup-run-client-tests cleanup-run-client-tests test-client build-linux build-osx build-windows internal-test-web-client vet run-server-for-web-client-tests

ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

IS_CI ?= false
# Build Flags
BUILD_NUMBER ?= $(BUILD_NUMBER:)
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse HEAD)
# If we don't set the build number it defaults to dev
ifeq ($(BUILD_NUMBER),)
	BUILD_NUMBER := dev
endif
BUILD_ENTERPRISE_DIR ?= ../enterprise
BUILD_ENTERPRISE ?= true
BUILD_ENTERPRISE_READY = false
BUILD_TYPE_NAME = team
BUILD_HASH_ENTERPRISE = none
LDAP_DATA ?= test
ifneq ($(wildcard $(BUILD_ENTERPRISE_DIR)/.),)
	ifeq ($(BUILD_ENTERPRISE),true)
		BUILD_ENTERPRISE_READY = true
		BUILD_TYPE_NAME = enterprise
		BUILD_HASH_ENTERPRISE = $(shell cd $(BUILD_ENTERPRISE_DIR) && git rev-parse HEAD)
	else
		BUILD_ENTERPRISE_READY = false
		BUILD_TYPE_NAME = team
	endif
else
	BUILD_ENTERPRISE_READY = false
	BUILD_TYPE_NAME = team
endif
BUILD_WEBAPP_DIR ?= ../xenia-webapp
BUILD_CLIENT = false
BUILD_HASH_CLIENT = independant
ifneq ($(wildcard $(BUILD_WEBAPP_DIR)/.),)
	ifeq ($(BUILD_CLIENT),true)
		BUILD_CLIENT = true
		BUILD_HASH_CLIENT = $(shell cd $(BUILD_WEBAPP_DIR) && git rev-parse HEAD)
	else
		BUILD_CLIENT = false
	endif
else
	BUILD_CLIENT = false
endif

# Golang Flags
GOPATH ?= $(shell go env GOPATH)
GOFLAGS ?= $(GOFLAGS:)
GO=go
DELVE=dlv
LDFLAGS += -X "github.com/xzl8028/xenia-server/model.BuildNumber=$(BUILD_NUMBER)"
LDFLAGS += -X "github.com/xzl8028/xenia-server/model.BuildDate=$(BUILD_DATE)"
LDFLAGS += -X "github.com/xzl8028/xenia-server/model.BuildHash=$(BUILD_HASH)"
LDFLAGS += -X "github.com/xzl8028/xenia-server/model.BuildHashEnterprise=$(BUILD_HASH_ENTERPRISE)"
LDFLAGS += -X "github.com/xzl8028/xenia-server/model.BuildEnterpriseReady=$(BUILD_ENTERPRISE_READY)"

# GOOS/GOARCH of the build host, used to determine whether we're cross-compiling or not
BUILDER_GOOS_GOARCH="$(shell $(GO) env GOOS)_$(shell $(GO) env GOARCH)"

PLATFORM_FILES="./cmd/xenia/main.go"

# Output paths
DIST_ROOT=dist
DIST_PATH=$(DIST_ROOT)/xenia

# Tests
TESTS=.

TESTFLAGS ?= -short
TESTFLAGSEE ?= -short

# Packages lists
TE_PACKAGES=$(shell go list ./...|grep -v plugin_tests)

# Plugins Packages
PLUGIN_PACKAGES=xenia-plugin-zoom-v1.0.7
PLUGIN_PACKAGES += xenia-plugin-autolink-v1.0.0
PLUGIN_PACKAGES += xenia-plugin-nps-v1.0.0
PLUGIN_PACKAGES += xenia-plugin-custom-attributes-v1.0.0
PLUGIN_PACKAGES += xenia-plugin-github-v0.10.2
PLUGIN_PACKAGES += xenia-plugin-welcomebot-v1.0.0
PLUGIN_PACKAGES += xenia-plugin-aws-SNS-v1.0.0
PLUGIN_PACKAGES += xenia-plugin-jira-v2.0.6

# Prepares the enterprise build if exists. The IGNORE stuff is a hack to get the Makefile to execute the commands outside a target
ifeq ($(BUILD_ENTERPRISE_READY),true)
	IGNORE:=$(shell echo Enterprise build selected, preparing)
	IGNORE:=$(shell rm -f imports/imports.go)
	IGNORE:=$(shell cp $(BUILD_ENTERPRISE_DIR)/imports/imports.go imports/)
	IGNORE:=$(shell rm -f enterprise)
	IGNORE:=$(shell ln -s $(BUILD_ENTERPRISE_DIR) enterprise)
else
	IGNORE:=$(shell rm -f imports/imports.go)
endif

EE_PACKAGES=$(shell go list ./enterprise/...)

ifeq ($(BUILD_ENTERPRISE_READY),true)
ALL_PACKAGES=$(TE_PACKAGES) $(EE_PACKAGES)
else
ALL_PACKAGES=$(TE_PACKAGES)
endif


all: run ## Alias for 'run'.

include build/*.mk

start-docker: ## Starts the docker containers for local development.
ifeq ($(IS_CI),false)
	@echo Starting docker containers

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-mysql$$ | wc -l) -eq 0 ]; then \
		echo starting xenia-mysql; \
		docker run --name xenia-mysql -p 3306:3306 \
			-e MYSQL_ROOT_PASSWORD=mostest \
			-e MYSQL_USER=mmuser \
			-e MYSQL_PASSWORD=mostest \
			-e MYSQL_DATABASE=xenia_test \
			-d mysql:5.7 > /dev/null; \
	elif [ $(shell docker ps --no-trunc --quiet --filter name=^/xenia-mysql$$ | wc -l) -eq 0 ]; then \
		echo restarting xenia-mysql; \
		docker start xenia-mysql > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-postgres$$ | wc -l) -eq 0 ]; then \
		echo starting xenia-postgres; \
		docker run --name xenia-postgres -p 5432:5432 \
			-e POSTGRES_USER=mmuser \
			-e POSTGRES_PASSWORD=mostest \
			-e POSTGRES_DB=xenia_test \
			-d postgres:9.4 > /dev/null; \
	elif [ $(shell docker ps --no-trunc --quiet --filter name=^/xenia-postgres$$ | wc -l) -eq 0 ]; then \
		echo restarting xenia-postgres; \
		docker start xenia-postgres > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-inbucket$$ | wc -l) -eq 0 ]; then \
		echo starting xenia-inbucket; \
		docker run --name xenia-inbucket -p 9000:10080 -p 2500:10025 -d jhillyerd/inbucket:release-1.2.0 > /dev/null; \
	elif [ $(shell docker ps --no-trunc --quiet --filter name=^/xenia-inbucket$$ | wc -l) -eq 0 ]; then \
		echo restarting xenia-inbucket; \
		docker start xenia-inbucket > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-minio$$ | wc -l) -eq 0 ]; then \
		echo starting xenia-minio; \
		docker run --name xenia-minio -p 9001:9000 -e "MINIO_ACCESS_KEY=minioaccesskey" \
		-e "MINIO_SSE_MASTER_KEY=my-minio-key:6368616e676520746869732070617373776f726420746f206120736563726574" \
		-e "MINIO_SECRET_KEY=miniosecretkey" -d minio/minio:RELEASE.2019-04-23T23-50-36Z server /data > /dev/null; \
		docker exec -it xenia-minio /bin/sh -c "mkdir -p /data/xenia-test" > /dev/null; \
	elif [ $(shell docker ps --no-trunc --quiet --filter name=^/xenia-minio$$ | wc -l) -eq 0 ]; then \
		echo restarting xenia-minio; \
		docker start xenia-minio > /dev/null; \
	fi

ifeq ($(BUILD_ENTERPRISE_READY),true)
	@echo Ldap test user test.one
	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-openldap$$ | wc -l) -eq 0 ]; then \
		echo starting xenia-openldap; \
		docker run --name xenia-openldap -p 389:389 -p 636:636 \
			-e LDAP_TLS_VERIFY_CLIENT="never" \
			-e LDAP_ORGANISATION="Xenia Test" \
			-e LDAP_DOMAIN="mm.test.com" \
			-e LDAP_ADMIN_PASSWORD="mostest" \
			-d osixia/openldap:1.2.2 > /dev/null;\
		sleep 10; \
		docker cp tests/test-data.ldif xenia-openldap:/test-data.ldif;\
		docker cp tests/qa-data.ldif xenia-openldap:/qa-data.ldif;\
		docker exec -ti xenia-openldap bash -c 'ldapadd -x -D "cn=admin,dc=mm,dc=test,dc=com" -w mostest -f /$(LDAP_DATA)-data.ldif';\
	elif [ $(shell docker ps | grep -ci xenia-openldap) -eq 0 ]; then \
		echo restarting xenia-openldap; \
		docker start xenia-openldap > /dev/null; \
		sleep 10; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-elasticsearch$$ | wc -l) -eq 0 ]; then \
		echo starting xenia-elasticsearch; \
		docker run --name xenia-elasticsearch -p 9200:9200 -e "http.host=0.0.0.0" -e "transport.host=127.0.0.1" -e "ES_JAVA_OPTS=-Xms250m -Xmx250m" -d xenia/xenia-elasticsearch-docker:6.5.1 > /dev/null; \
	elif [ $(shell docker ps --no-trunc --quiet --filter name=^/xenia-elasticsearch$$ | wc -l) -eq 0 ]; then \
		echo restarting xenia-elasticsearch; \
		docker start xenia-elasticsearch> /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-redis$$ | wc -l) -eq 0 ]; then \
		echo starting xenia-redis; \
		docker run --name xenia-redis -p 6379:6379 -d redis > /dev/null; \
	elif [ $(shell docker ps --no-trunc --quiet --filter name=^/xenia-redis$$ | wc -l) -eq 0 ]; then \
		echo restarting xenia-redis; \
		docker start xenia-redis > /dev/null; \
	fi
endif
else
	@echo CI Build: skipping docker start
endif

stop-docker: ## Stops the docker containers for local development.
	@echo Stopping docker containers

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-mysql$$ | wc -l) -eq 1 ]; then \
		echo stopping xenia-mysql; \
		docker stop xenia-mysql > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-mysql-unittest$$ | wc -l) -eq 1 ]; then \
		echo stopping xenia-mysql-unittest; \
		docker stop xenia-mysql-unittest > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-postgres$$ | wc -l) -eq 1 ]; then \
		echo stopping xenia-postgres; \
		docker stop xenia-postgres > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-postgres-unittest$$ | wc -l) -eq 1 ]; then \
		echo stopping xenia-postgres-unittest; \
		docker stop xenia-postgres-unittest > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-openldap$$ | wc -l) -eq 1 ]; then \
		echo stopping xenia-openldap; \
		docker stop xenia-openldap > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-inbucket$$ | wc -l) -eq 1 ]; then \
		echo stopping xenia-inbucket; \
		docker stop xenia-inbucket > /dev/null; \
	fi

		@if [ $(shell docker ps -a | grep -ci xenia-minio) -eq 1 ]; then \
		echo stopping xenia-minio; \
		docker stop xenia-minio > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-elasticsearch$$ | wc -l) -eq 1 ]; then \
		echo stopping xenia-elasticsearch; \
		docker stop xenia-elasticsearch > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-redis$$ | wc -l) -eq 1 ]; then \
		echo stopping xenia-redis; \
		docker stop xenia-redis > /dev/null; \
	fi

clean-docker: ## Deletes the docker containers for local development.
	@echo Removing docker containers

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-mysql$$ | wc -l) -eq 1 ]; then \
		echo removing xenia-mysql; \
		docker stop xenia-mysql > /dev/null; \
		docker rm -v xenia-mysql > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-mysql-unittest$$ | wc -l) -eq 1 ]; then \
		echo removing xenia-mysql-unittest; \
		docker stop xenia-mysql-unittest > /dev/null; \
		docker rm -v xenia-mysql-unittest > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-postgres$$ | wc -l) -eq 1 ]; then \
		echo removing xenia-postgres; \
		docker stop xenia-postgres > /dev/null; \
		docker rm -v xenia-postgres > /dev/null; \
	fi

	@if [ $(shell docker ps -a --no-trunc --quiet --filter name=^/xenia-postgres-unittest$$ | wc -l) -eq 1 ]; then \
		echo removing xenia-postgres-unittest; \
		docker stop xenia-postgres-unittest > /dev/null; \
		docker rm -v xenia-postgres-unittest > /dev/null; \
	fi

	@if [ $(shell docker ps -a | grep -ci xenia-openldap) -eq 1 ]; then \
		echo removing xenia-openldap; \
		docker stop xenia-openldap > /dev/null; \
		docker rm -v xenia-openldap > /dev/null; \
	fi

	@if [ $(shell docker ps -a | grep -ci xenia-inbucket) -eq 1 ]; then \
		echo removing xenia-inbucket; \
		docker stop xenia-inbucket > /dev/null; \
		docker rm -v xenia-inbucket > /dev/null; \
	fi

	@if [ $(shell docker ps -a | grep -ci xenia-minio) -eq 1 ]; then \
		echo removing xenia-minio; \
		docker stop xenia-minio > /dev/null; \
		docker rm -v xenia-minio > /dev/null; \
	fi

	@if [ $(shell docker ps -a | grep -ci xenia-elasticsearch) -eq 1 ]; then \
		echo removing xenia-elasticsearch; \
		docker stop xenia-elasticsearch > /dev/null; \
		docker rm -v xenia-elasticsearch > /dev/null; \
	fi

govet: ## Runs govet against all packages.
	@echo Running GOVET
	$(GO) get golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	$(GO) vet $(GOFLAGS) $(ALL_PACKAGES) || exit 1
	$(GO) vet -vettool=$(GOPATH)/bin/shadow $(GOFLAGS) $(ALL_PACKAGES) || exit 1

gofmt: ## Runs gofmt against all packages.
	@echo Running GOFMT

	@for package in $(TE_PACKAGES) $(EE_PACKAGES); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failure"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success"; \

megacheck: ## Run megacheck on codebasis
	env GO111MODULE=off go get -u honnef.co/go/tools/cmd/megacheck
	$(GOPATH)/bin/megacheck $(TE_PACKAGES)

ifeq ($(BUILD_ENTERPRISE_READY),true)
	$(GOPATH)/bin/megacheck $(EE_PACKAGES) || exit 1
endif

i18n-extract: ## Extract strings for translation from the source code
	env GO111MODULE=off go get -u github.com/xzl8028/xenia-utilities/mmgotool
	$(GOPATH)/bin/mmgotool i18n extract

store-mocks: ## Creates mock files.
	env GO111MODULE=off go get -u github.com/vektra/mockery/...
	$(GOPATH)/bin/mockery -dir store -all -output store/storetest/mocks -note 'Regenerate this file using `make store-mocks`.'

filesstore-mocks: ## Creates mock files.
	env GO111MODULE=off go get -u github.com/vektra/mockery/...
	$(GOPATH)/bin/mockery -dir services/filesstore -all -output services/filesstore/mocks -note 'Regenerate this file using `make filesstore-mocks`.'

ldap-mocks: ## Creates mock files for ldap.
	env GO111MODULE=off go get -u github.com/vektra/mockery/...
	$(GOPATH)/bin/mockery -dir enterprise/ldap -all -output enterprise/ldap/mocks -note 'Regenerate this file using `make ldap-mocks`.'

plugin-mocks: ## Creates mock files for plugins.
	env GO111MODULE=off go get -u github.com/vektra/mockery/...
	$(GOPATH)/bin/mockery -dir plugin -name API -output plugin/plugintest -outpkg plugintest -case underscore -note 'Regenerate this file using `make plugin-mocks`.'
	$(GOPATH)/bin/mockery -dir plugin -name Hooks -output plugin/plugintest -outpkg plugintest -case underscore -note 'Regenerate this file using `make plugin-mocks`.'
	$(GOPATH)/bin/mockery -dir plugin -name Helpers -output plugin/plugintest -outpkg plugintest -case underscore -note 'Regenerate this file using `make plugin-mocks`.'

pluginapi: ## Generates api and hooks glue code for plugins
	go generate ./plugin

check-licenses: ## Checks license status.
	./scripts/license-check.sh $(TE_PACKAGES) $(EE_PACKAGES)

check-prereqs: ## Checks prerequisite software status.
	./scripts/prereq-check.sh

check-style: govet gofmt check-licenses ## Runs govet and gofmt against all packages.

test-te-race: ## Checks for race conditions in the team edition.
	@echo Testing TE race conditions

	@echo "Packages to test: "$(TE_PACKAGES)

	@for package in $(TE_PACKAGES); do \
		echo "Testing "$$package; \
		$(GO) test $(GOFLAGS) -race -run=$(TESTS) -test.timeout=4000s $$package || exit 1; \
	done

test-ee-race: ## Checks for race conditions in the enterprise edition.
	@echo Testing EE race conditions

ifeq ($(BUILD_ENTERPRISE_READY),true)
	@echo "Packages to test: "$(EE_PACKAGES)

	for package in $(EE_PACKAGES); do \
		echo "Testing "$$package; \
		$(GO) test $(GOFLAGS) -race -run=$(TESTS) -c $$package; \
		if [ -f $$(basename $$package).test ]; then \
			echo "Testing "$$package; \
			./$$(basename $$package).test -test.timeout=2000s || exit 1; \
			rm -r $$(basename $$package).test; \
		fi; \
	done

	rm -f config/*.crt
	rm -f config/*.key
endif

test-server-race: test-te-race test-ee-race ## Checks for race conditions.
	find . -type d -name data -not -path './vendor/*' | xargs rm -rf

do-cover-file: ## Creates the test coverage report file.
	@echo "mode: count" > cover.out

go-junit-report:
	env GO111MODULE=off go get -u github.com/jstemmer/go-junit-report

test-compile:
	@echo COMPILE TESTS

	for package in $(TE_PACKAGES) $(EE_PACKAGES); do \
		$(GO) test $(GOFLAGS) -c $$package; \
	done

test-db-migration: start-docker
	./scripts/mysql-migration-test.sh
	./scripts/psql-migration-test.sh

test-server: start-docker go-junit-report do-cover-file ## Runs tests.
ifeq ($(BUILD_ENTERPRISE_READY),true)
	@echo Running all tests
else
	@echo Running only TE tests
endif
	./scripts/test.sh "$(GO)" "$(GOFLAGS)" "$(ALL_PACKAGES)" "$(TESTS)" "$(TESTFLAGS)"

internal-test-web-client: ## Runs web client tests.
	$(GO) run $(GOFLAGS) $(PLATFORM_FILES) test web_client_tests

run-server-for-web-client-tests: ## Tests the server for web client.
	$(GO) run $(GOFLAGS) $(PLATFORM_FILES) test web_client_tests_server

test-client: ## Test client app.
	@echo Running client tests

	cd $(BUILD_WEBAPP_DIR) && $(MAKE) test

test: test-server test-client ## Runs all checks and tests below (except race detection and postgres).

cover: ## Runs the golang coverage tool. You must run the unit tests first.
	@echo Opening coverage info in browser. If this failed run make test first

	$(GO) tool cover -html=cover.out
	$(GO) tool cover -html=ecover.out

test-data: start-docker ## Add test data to the local instance.
	$(GO) run $(GOFLAGS) -ldflags '$(LDFLAGS)' $(PLATFORM_FILES) config set TeamSettings.MaxUsersPerTeam 100
	$(GO) run $(GOFLAGS) -ldflags '$(LDFLAGS)' $(PLATFORM_FILES) sampledata -w 4 -u 60 

	@echo You may need to restart the Xenia server before using the following
	@echo ========================================================================
	@echo Login with a system admin account username=sysadmin password=sysadmin
	@echo Login with a regular account username=user-1 password=user-1
	@echo ========================================================================

run-server: start-docker ## Starts the server.
	@echo Running xenia for development

	mkdir -p $(BUILD_WEBAPP_DIR)/dist/files
	$(GO) run $(GOFLAGS) -ldflags '$(LDFLAGS)' $(PLATFORM_FILES) --disableconfigwatch | \
	    $(GO) run $(GOFLAGS) -ldflags '$(LDFLAGS)' $(PLATFORM_FILES) logs --logrus &

debug-server: start-docker
	mkdir -p $(BUILD_WEBAPP_DIR)/dist/files
	$(DELVE) debug $(PLATFORM_FILES) --build-flags="-ldflags '\
		-X github.com/xzl8028/xenia-server/model.BuildNumber=$(BUILD_NUMBER)\
		-X \"github.com/xzl8028/xenia-server/model.BuildDate=$(BUILD_DATE)\"\
		-X github.com/xzl8028/xenia-server/model.BuildHash=$(BUILD_HASH)\
		-X github.com/xzl8028/xenia-server/model.BuildHashEnterprise=$(BUILD_HASH_ENTERPRISE)\
		-X github.com/xzl8028/xenia-server/model.BuildEnterpriseReady=$(BUILD_ENTERPRISE_READY)'"

run-cli: start-docker ## Runs CLI.
	@echo Running xenia for development
	@echo Example should be like 'make ARGS="-version" run-cli'

	$(GO) run $(GOFLAGS) -ldflags '$(LDFLAGS)' $(PLATFORM_FILES) ${ARGS}

run-client: ## Runs the webapp.
	@echo Running xenia client for development

	ln -nfs $(BUILD_WEBAPP_DIR)/dist client
	cd $(BUILD_WEBAPP_DIR) && $(MAKE) run

run-client-fullmap: ## Legacy alias to run-client
	@echo Running xenia client for development

	cd $(BUILD_WEBAPP_DIR) && $(MAKE) run

run: check-prereqs run-server run-client ## Runs the server and webapp.

run-fullmap: run-server run-client ## Legacy alias to run

stop-server: ## Stops the server.
	@echo Stopping xenia

ifeq ($(BUILDER_GOOS_GOARCH),"windows_amd64")
	wmic process where "Caption='go.exe' and CommandLine like '%go.exe run%'" call terminate
	wmic process where "Caption='xenia.exe' and CommandLine like '%go-build%'" call terminate
else
	@for PID in $$(ps -ef | grep "[g]o run" | awk '{ print $$2 }'); do \
		echo stopping go $$PID; \
		kill $$PID; \
	done
	@for PID in $$(ps -ef | grep "[g]o-build" | awk '{ print $$2 }'); do \
		echo stopping xenia $$PID; \
		kill $$PID; \
	done
endif

stop-client: ## Stops the webapp.
	@echo Stopping xenia client

	cd $(BUILD_WEBAPP_DIR) && $(MAKE) stop

stop: stop-server stop-client ## Stops server and client.

restart: restart-server restart-client ## Restarts the server and webapp.

restart-server: | stop-server run-server ## Restarts the xenia server to pick up development change.

restart-client: | stop-client run-client ## Restarts the webapp.

run-job-server: ## Runs the background job server.
	@echo Running job server for development
	$(GO) run $(GOFLAGS) -ldflags '$(LDFLAGS)' $(PLATFORM_FILES) jobserver --disableconfigwatch &

config-ldap: ## Configures LDAP.
	@echo Setting up configuration for local LDAP

	@sed -i'' -e 's|"LdapServer": ".*"|"LdapServer": "dockerhost"|g' config/config.json
	@sed -i'' -e 's|"BaseDN": ".*"|"BaseDN": "dc=mm,dc=test,dc=com"|g' config/config.json
	@sed -i'' -e 's|"BindUsername": ".*"|"BindUsername": "cn=admin,dc=mm,dc=test,dc=com"|g' config/config.json
	@sed -i'' -e 's|"BindPassword": ".*"|"BindPassword": "mostest"|g' config/config.json
	@sed -i'' -e 's|"FirstNameAttribute": ".*"|"FirstNameAttribute": "cn"|g' config/config.json
	@sed -i'' -e 's|"LastNameAttribute": ".*"|"LastNameAttribute": "sn"|g' config/config.json
	@sed -i'' -e 's|"NicknameAttribute": ".*"|"NicknameAttribute": "cn"|g' config/config.json
	@sed -i'' -e 's|"EmailAttribute": ".*"|"EmailAttribute": "mail"|g' config/config.json
	@sed -i'' -e 's|"UsernameAttribute": ".*"|"UsernameAttribute": "uid"|g' config/config.json
	@sed -i'' -e 's|"IdAttribute": ".*"|"IdAttribute": "uid"|g' config/config.json
	@sed -i'' -e 's|"LoginIdAttribute": ".*"|"LoginIdAttribute": "uid"|g' config/config.json
	@sed -i'' -e 's|"GroupDisplayNameAttribute": ".*"|"GroupDisplayNameAttribute": "cn"|g' config/config.json
	@sed -i'' -e 's|"GroupIdAttribute": ".*"|"GroupIdAttribute": "entryUUID"|g' config/config.json

config-reset: ## Resets the config/config.json file to the default.
	@echo Resetting configuration to default
	rm -f config/config.json
	OUTPUT_CONFIG=$(PWD)/config/config.json go generate ./config

clean: stop-docker ## Clean up everything except persistant server data.
	@echo Cleaning

	rm -Rf $(DIST_ROOT)
	go clean $(GOFLAGS) -i ./...

	cd $(BUILD_WEBAPP_DIR) && $(MAKE) clean

	find . -type d -name data -not -path './vendor/*' | xargs rm -rf
	rm -rf logs

	rm -f xenia.log
	rm -f xenia.log.jsonl
	rm -f npm-debug.log
	rm -f .prepare-go
	rm -f enterprise
	rm -f cover.out
	rm -f ecover.out
	rm -f *.out
	rm -f *.test
	rm -f imports/imports.go
	rm -f cmd/platform/cprofile*.out
	rm -f cmd/xenia/cprofile*.out

nuke: clean clean-docker ## Clean plus removes persistent server data.
	@echo BOOM

	rm -rf data

setup-mac: ## Adds macOS hosts entries for Docker.
	echo $$(boot2docker ip 2> /dev/null) dockerhost | sudo tee -a /etc/hosts

update-dependencies: ## Uses go get -u to update all the dependencies while holding back any that require it.
	@echo Updating Dependencies

	# Update all dependencies (does not update across major versions)
	go get -u

	# Keep back because of breaking API changes
	go get -u github.com/segmentio/analytics-go@2.1.1

	# Tidy up
	go mod tidy

	# Copy everything to vendor directory
	go mod vendor


todo: ## Display TODO and FIXME items in the source code.
	@! ag --ignore Makefile --ignore-dir vendor --ignore-dir runtime TODO
	@! ag --ignore Makefile --ignore-dir vendor --ignore-dir runtime XXX
	@! ag --ignore Makefile --ignore-dir vendor --ignore-dir runtime FIXME
	@! ag --ignore Makefile --ignore-dir vendor --ignore-dir runtime "FIX ME"
ifeq ($(BUILD_ENTERPRISE_READY),true)
	@! ag --ignore Makefile --ignore-dir vendor --ignore-dir runtime TODO enterprise/
	@! ag --ignore Makefile --ignore-dir vendor --ignore-dir runtime XXX enterprise/
	@! ag --ignore Makefile --ignore-dir vendor --ignore-dir runtime FIXME enterprise/
	@! ag --ignore Makefile --ignore-dir vendor --ignore-dir runtime "FIX ME" enterprise/
endif

## Help documentatin à la https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' ./Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
