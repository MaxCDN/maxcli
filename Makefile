GOROOT=$(shell go env GOROOT)
PLATFORMS=darwin/386 \
			darwin/amd64 \
			freebsd/386 \
			freebsd/amd64 \
			freebsd/arm \
			linux/386 \
			linux/amd64 \
			linux/arm \
			windows/386 \
			windows/amd64

# todo: find a cleaner way
BUILD_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./build/')
FORMAT_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./format/')
DEPLOY_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./deploy/')
INSTALL_TOOLS=$(shell find . -maxdepth 1 -type d -name "max*" | sed 's/\./install/')

# tests
###
test: format .PHONY
	./_test/shunt.sh --verbose ./_test/*_test.sh

test/%: .PHONY
	# Running Test: $*
	./_test/shunt.sh --verbose ./_test/$*_test.sh

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
	bash build-all.bash $*

build/all: .PHONY
	# exec tools/build-all.sh
	bash build-all.bash

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
	@find . -maxdepth 2 -type d -name builds -exec rm -vr {} \; || true

# setup and sub-tasks setup go for crosscompiling
###
setup: $(PLATFORMS) .PHONY

darwin/386:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=darwin GOARCH=386 ./make.bash --no-clean 2>&1'

darwin/amd64:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=darwin GOARCH=amd64 ./make.bash --no-clean 2>&1'

freebsd/386:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=freebsd GOARCH=386 ./make.bash --no-clean 2>&1'

freebsd/amd64:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=freebsd GOARCH=amd64 ./make.bash --no-clean 2>&1'

freebsd/arm:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=freebsd GOARCH=arm ./make.bash --no-clean 2>&1'

linux/386:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=linux GOARCH=386 ./make.bash --no-clean 2>&1'

linux/amd64:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=linux GOARCH=amd64 ./make.bash --no-clean 2>&1'

linux/arm:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=linux GOARCH=arm ./make.bash --no-clean 2>&1'

windows/386:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=windows GOARCH=386 ./make.bash --no-clean 2>&1'

windows/amd64:
	sudo bash -c 'cd $(GOROOT)/src && GOOS=windows GOARCH=amd64 ./make.bash --no-clean 2>&1'

# deploy
##
deploy: $(DEPLOY_TOOLS) .PHONY

deploy/%: .PHONY
	@which "aws" > /dev/null || (echo "ERROR: install 'awscli' tools." && exit 1)
	aws s3 cp $*/builds/ s3://maxcli/$*/ \
		--recursive \
		--acl public-read

.PHONY:
