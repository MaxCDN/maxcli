GOROOT=$(shell go env GOROOT)

# todo: find a cleaner way
BUILD_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./build/')
FORMAT_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./format/')
DEPLOY_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./deploy/')
INSTALL_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./install/')
GET_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./get/')

# tests
###
test: format .PHONY
	./_test/shunt.sh --verbose ./_test/*_test.sh

test/%: .PHONY
	# Running Test: $*
	./_test/shunt.sh --verbose ./_test/$*_test.sh

# get
###
get: $(GET_TOOLS) .PHONY

get/%: .PHONY
	cd $*; dep ensure

# format
###
format: $(FORMAT_TOOLS) .PHONY

format/%: .PHONY
	# Formatting: $*
	@cd $*; go fmt

# build tools
###
build: $(BUILD_TOOLS) .PHONY

install: $(INSTALL_TOOLS) .PHONY

build/all/%: .PHONY
	make format/$*
	env GOROOT=$(GOROOT) bash scripts/build-all.bash $*

build/all: .PHONY
	env GOROOT=$(GOROOT) bash scripts/build-all.bash

build/%: .PHONY
	# Building: $*
	make format/$*
	cd $*; go build $*.go

install/%: .PHONY
	# Installing: $*
	@test "$(GOBIN)" || (echo "error: GOBIN must be set"; exit 1)
	@test -x $*/$* || (echo "error: run make build[/{{tool}}]"; exit 1)
	mv $*/$* $(GOBIN)

clean:
	# remove previous tools builds
	@test -f maxcurl/maxcurl && rm -v maxcurl/maxcurl || true
	@test -f maxpurge/maxpurge && rm -v maxpurge/maxpurge || true
	@test -f maxreport/maxreport && rm -v maxreport/maxreport || true
	@test -f maxtail/maxtail && rm -v maxtail/maxtail || true
	@rm -rf _builds

.PHONY:
