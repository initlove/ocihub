PREFIX ?= $(DESTDIR)/usr
BINDIR ?= $(DESTDIR)/usr/bin

BUILDTAGS=
COMMIT=$(shell git rev-parse HEAD 2> /dev/null || true)

all: ocihub

ocihub:
	go build -tags "$(BUILDTAGS)" -ldflags "-X main.gitCommit=${COMMIT}" -o ocihub

clean:
	rm -f ocihub

.PHONY: test .gofmt .govet .golint

PACKAGES = $(shell go list ./... | grep -v vendor)
test: .gofmt .govet .golint .gotest

.gofmt:
	OUT=$$(go fmt $(PACKAGES)); if test -n "$${OUT}"; then echo "$${OUT}" && exit 1; fi

.govet:
	go vet -x $(PACKAGES)

.golint:
	golint -set_exit_status $(PACKAGES)

.gotest:
	go test $(PACKAGES)
