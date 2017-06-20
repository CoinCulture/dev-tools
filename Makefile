PACKAGES=$(shell go list ./... | grep -v '/vendor/')

all: get_vendor_deps install test

build:
	go build ./cmd/...

install:
	go install ./cmd/...

test:
	go test $(PACKAGES)

get_vendor_deps:
	go get -u -v github.com/Masterminds/glide
	glide install

.PHONY: all build install test get_vendor_deps
