GO ?= $(shell command -v go 2> /dev/null)
NPM ?= $(shell command -v npm 2> /dev/null)
CURL ?= $(shell command -v curl 2> /dev/null)
GOARCH := amd64
GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
GOPATH ?= $(shell go env GOPATH)
GO_TEST_FLAGS ?= -race
GO_BUILD_FLAGS ?=
MM_UTILITIES_DIR ?= ../mattermost-utilities

export GO111MODULE=on

MANIFEST_FILE ?= plugin.json

# You can include assets this directory into the bundle. This can be e.g. used to include profile pictures.
ASSETS_DIR ?= assets

PLUGIN_ID=$(shell echo `grep '"'"id"'"\s*:\s*"' $(MANIFEST_FILE) | head -1 | cut -d'"' -f4`)
PLUGIN_VERSION=v$(shell echo `grep '"'"version"'"\s*:\s*"' $(MANIFEST_FILE) | head -1 | cut -d'"' -f4`)
BUNDLE_NAME ?= mattermost-plugin-$(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz

HAS_WEBAPP=$(shell if [ "$(shell grep -E '[\^"]webapp["][ ]*[:]' $(MANIFEST_FILE)  | wc -l)" -gt "0" ]; then echo "true"; fi)
HAS_SERVER=$(shell if [ "$(shell grep -E '[\^"]server["][ ]*[:]' $(MANIFEST_FILE)  | wc -l)" -gt "0" ]; then echo "true"; fi)

TMPFILEGOLINT=golint.tmp

BLACK=`tput setaf 0`
RED=`tput setaf 1`
GREEN=`tput setaf 2`
YELLOW=`tput setaf 3`
BLUE=`tput setaf 4`
MAGENTA=`tput setaf 5`
CYAN=`tput setaf 6`
WHITE=`tput setaf 7`

BOLD=`tput bold`
INVERSE=`tput rev`
RESET=`tput sgr0`

## Runs 'make all' by default
.PHONY: default
default: all

## Checks the code style, tests, builds and bundles the plugin.
.PHONY: all
all: clean check-style test dist

## Propagates plugin manifest information into the server/ and webapp/ folders as required.
.PHONY: apply
apply:
ifneq ($(HAS_WEBAPP),)
	@mkdir -p webapp/src/constants
	@echo "export const PLUGIN_NAME = '`echo $(PLUGIN_ID)`';" > webapp/src/constants/manifest.js
endif
ifneq ($(HAS_SERVER),)
	@echo "package config\n\nconst (\n\tPluginName = \""`echo $(PLUGIN_ID)`"\"\n)" > server/config/manifest.go
endif

## Triggers a Release through CircleCI by pushing a GitHub Tag
.PHONY: trigger-release
trigger-release:
	@if [ $$(git status --porcelain | wc -l) != "0" -o $$(git rev-list HEAD@{upstream}..HEAD | wc -l) != "0" ]; \
		then echo ${RED}"local repo is not clean"${RESET}; exit 1; fi;
	@echo ${BOLD}"Creating a tag to trigger circleci build-and-release job\n"${RESET}
	git tag $(PLUGIN_VERSION)
	git push origin $(PLUGIN_VERSION)

## Checks for style guide compliance in webapp and server.
.PHONY: check-style
check-style: webapp/.npminstall check-style-webapp gofmt govet golint

## Checks for style guide compliance in the webapp.
.PHONY: check-style-webapp
check-style-webapp:
ifneq ($(HAS_WEBAPP),)
	@echo ${BOLD}Running ESLINT${RESET}
	@cd webapp && npm run lint
	@echo ${GREEN}"eslint success\n"${RESET}
endif

## Runs gofmt against all packages.
.PHONY: gofmt
gofmt:
ifneq ($(HAS_SERVER),)
	@echo ${BOLD}Running GOFMT${RESET}
	@for package in $$(go list ./server/...); do \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo ${RED}"gofmt failure\n"${RESET}; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo ${GREEN}"gofmt success\n"${RESET}
endif

## Runs govet against all packages.
.PHONY: govet
govet:
ifneq ($(HAS_SERVER),)
	@echo ${BOLD}Running GOVET${RESET}
	@# Workaround because you can't install binaries without adding them to go.mod
	env GO111MODULE=off $(GO) get golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	cd server && $(GO) vet ./...
	cd server && $(GO) vet -vettool=$(GOPATH)/bin/shadow ./...
	@echo ${GREEN}"govet success\n"${RESET}
endif

## Runs golint against all packages.
.PHONY: golint
golint:
ifneq ($(HAS_SERVER),)
	@echo ${BOLD}Running GOLINT${RESET}
	env GO111MODULE=off $(GO) get golang.org/x/lint/golint
	$(eval PKGS := $(shell go list ./... | grep -v /vendor/))
	@touch $(TMPFILEGOLINT)
	-@for pkg in $(PKGS) ; do \
		echo `$(GOPATH)/bin/golint $$pkg | grep -v "have comment" | grep -v "comment on exported" | grep -v "lint suggestions"` >> $(TMPFILEGOLINT) ; \
	done
	@grep -Ev "^$$" $(TMPFILEGOLINT) || true
	@if [ "$$(grep -Ev "^$$" $(TMPFILEGOLINT) | wc -l)" -gt "0" ]; then \
		rm -f $(TMPFILEGOLINT); echo ${RED}"golint failure\n"${RESET}; exit 1; else \
		rm -f $(TMPFILEGOLINT); echo ${GREEN}"golint success\n"${RESET}; \
	fi
endif

## Fixes webapp and server styles.
.PHONY: format
format: fix-webapp fix-go

## Fixes webapp styles.
.PHONY: fix-webapp
fix-webapp:
ifneq ($(HAS_WEBAPP),)
	@echo ${BOLD}Formatting js files${RESET}
	@cd webapp && npm run fix
	@echo "formatted js files\n"
endif

## Runs goimports against all packages.
.PHONY: fix-go
fix-go:
ifneq ($(HAS_SERVER),)
	env GO111MODULE=off $(GO) get golang.org/x/tools/cmd/goimports
	@echo ${BOLD}Formatting go files${RESET}
	@cd server
	@find ./ -type f -name "*.go" -not -path "./server/vendor/*" -exec goimports -w {} \;
	@echo "formatted go files\n"
endif

## Runs any lints and unit tests defined for the server and webapp, if they exist.
.PHONY: test
test: webapp/.npminstall
ifneq ($(HAS_SERVER),)
	$(GO) test -v $(GO_TEST_FLAGS) ./server/...
endif
ifneq ($(HAS_WEBAPP),)
	cd webapp && $(NPM) run fix && $(NPM) run test;
endif

## Creates a coverage report for the server code.
.PHONY: coverage
coverage: webapp/.npminstall
ifneq ($(HAS_SERVER),)
	$(GO) test $(GO_TEST_FLAGS) -coverprofile=server/coverage.txt ./server/...
	$(GO) tool cover -html=server/coverage.txt
endif

## Extract strings for translation from the source code.
.PHONY: i18n-extract
i18n-extract:
ifneq ($(HAS_WEBAPP),)
ifeq ($(HAS_MM_UTILITIES),)
	@echo "You must clone github.com/mattermost/mattermost-utilities repo in .. to use this command"
else
	cd $(MM_UTILITIES_DIR) && npm install && npm run babel && node mmjstool/build/index.js i18n extract-webapp --webapp-dir $(PWD)/webapp
endif
endif


## Ensures NPM dependencies are installed without having to run this all the time.
webapp/.npminstall:
ifneq ($(HAS_WEBAPP),)
	@echo ${BOLD}"Getting dependencies using npm\n"${RESET}
	cd webapp && npm install
	touch $@
	@echo "\n"
endif

## Builds the webapp, if it exists.
.PHONY: webapp
webapp: webapp/.npminstall
ifneq ($(HAS_WEBAPP),)
	cd webapp && $(NPM) run build;
endif

## Builds the webapp in debug mode, if it exists.
webapp-debug: webapp/.npminstall
ifneq ($(HAS_WEBAPP),)
	cd webapp && \
	$(NPM) run debug;
endif

## Gets the vendor dependencies.
vendor: go.sum
ifneq ($(HAS_SERVER),)
	$(GO) mod vendor
endif

## Builds the server, if it exists, including support for multiple architectures.
.PHONY: server
server:
ifneq ($(HAS_SERVER),)
	mkdir -p server/dist;
	cd server && env GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o dist/plugin-linux-amd64;
	cd server && env GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o dist/plugin-darwin-amd64;
	cd server && env GOOS=windows GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -o dist/plugin-windows-amd64.exe;
endif

## Generates a tar bundle of the plugin for install.
.PHONY: bundle
bundle:
	rm -rf dist/
	mkdir -p dist/$(PLUGIN_ID)
	cp $(MANIFEST_FILE) dist/$(PLUGIN_ID)/
ifneq ($(wildcard $(ASSETS_DIR)/.),)
	cp -r $(ASSETS_DIR) dist/$(PLUGIN_ID)/
endif
ifneq ($(HAS_PUBLIC),)
	cp -r public/ dist/$(PLUGIN_ID)/
endif
ifneq ($(HAS_SERVER),)
	mkdir -p dist/$(PLUGIN_ID)/server/dist;
	cp -r server/dist/* dist/$(PLUGIN_ID)/server/dist/;
endif
ifneq ($(HAS_WEBAPP),)
	mkdir -p dist/$(PLUGIN_ID)/webapp/dist;
	cp -r webapp/dist/* dist/$(PLUGIN_ID)/webapp/dist/;
endif
	cd dist && tar -cvzf $(BUNDLE_NAME) $(PLUGIN_ID)

	@echo plugin built at: dist/$(BUNDLE_NAME)

## Builds and bundles the plugin.
.PHONY: dist
dist: .distclean apply server webapp bundle

## Builds and bundles the plugin in debug mode (if it exists).
.PHONY: debug-dist
debug-dist: .distclean apply server webapp-debug bundle

## Removes all build artifacts.
.PHONY: .distclean
.distclean:
	@echo ${BOLD}"Cleaning dist files\n"${RESET}
	rm -fr dist/
	rm -fr server/dist
	rm -fr webapp/dist
	@echo "\n"

## Removes all dependencies and build-artifacts.
.PHONY: clean
clean: .distclean
	@echo ${BOLD}"Cleaning plugin\n"${RESET}
	rm -fr vendor
	rm -fr webapp/node_modules
	rm -fr webapp/.npminstall
	@echo "\n"

# Help documentation à la https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@cat Makefile | grep -v '\.PHONY' |  grep -v '\help:' | grep -B1 -E '^[a-zA-Z0-9_.-]+:.*' | sed -e "s/:.*//" | sed -e "s/^## //" |  grep -v '\-\-' | sed '1!G;h;$$!d' | awk 'NR%2{printf "\033[36m%-30s\033[0m",$$0;next;}1' | sort
