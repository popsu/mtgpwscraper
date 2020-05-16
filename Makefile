SHELL := /bin/sh

BINS = bin/macos/pwscraper-macos bin/win64/pwscraper-win64.exe bin/linux64/pwscraper-linux64
GOSRC = $(shell find . -name '*.go')
CGO_ENABLED = 0
ENVVARS = CGO_ENABLED=$(CGO_ENABLED)
LDFLAGS = -ldflags='-s -w'
BUILDFLAGS = -tags "osusergo netgo"
BUILDWINFLAGS = -tags "osusergo netgo"

PACKAGE = github.com/popsu/mtgpwscraper/cmd

.PHONY: build
build: $(BINS)

.PHONY: dev
dev: bin/linux64/pwscraper

.PHONY: clean
clean:
	- go clean .
	- rm -vrf bin/*

bin/macos/pwscraper-macos: $(GOSRC)
	$(ENVVARS) GOOS=darwin go build $(LDFLAGS) $(BUILDFLAGS) -o $@ $(PACKAGE)

bin/win64/pwscraper-win64.exe: $(GOSRC)
	$(ENVVARS) GOOS=windows CC=x86_64-w64-mingw32-gcc \
	go build $(LDFLAGS) $(BUILDWINFLAGS) -o $@ $(PACKAGE)

bin/linux64/pwscraper-linux64: $(GOSRC)
	$(ENVVARS) go build $(LDFLAGS) $(BUILDFLAGS) -o $@ $(PACKAGE)
