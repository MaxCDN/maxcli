GOPATH=$(shell pwd)/Godeps/workspace

# Run Tests
test: format
	go test

# Go Get Deps
get:
	go get -u -v github.com/garyburd/go-oauth/oauth
	go get -u -v gopkg.in/jmervine/GoT.v1

# Show Docs
docs: format
	@godoc -ex . | sed -e 's/func /\nfunc /g' | less # add a little spacing for readability

# Go Fmt Source
format:
	@gofmt -s -w -l $(shell find . -type f -name "*.go" | grep -v Godeps)

.PHONY: travis test docs format
