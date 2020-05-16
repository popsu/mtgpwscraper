SHELL := /bin/sh

BINS = bin/macos/pwscraper bin/win64/pwscraper.exe bin/linux64/pwscraper
GOSRC = $(shell find . -name '*.go')
CGO_ENABLED = 0
ENVVARS = CGO_ENABLED=$(CGO_ENABLED)
LDFLAGS = -ldflags='-s -w'
BUILDFLAGS = -tags "osusergo netgo"
BUILDWINFLAGS = -tags "osusergo netgo"

PACKAGE = github.com/popsu/mtgpwscraper/cmd

.PHONY: build
build: $(BINS)

.PHONY: clean
clean:
	- go clean .
	- rm -vrf bin/*

bin/macos/pwscraper: $(GOSRC)
	$(ENVVARS) GOOS=darwin go build $(LDFLAGS) $(BUILDFLAGS) -o $@ $(PACKAGE)

bin/win64/pwscraper.exe: $(GOSRC)
	$(ENVVARS) GOOS=windows CC=x86_64-w64-mingw32-gcc \
	go build $(LDFLAGS) $(BUILDWINFLAGS) -o $@ $(PACKAGE)

bin/linux64/pwscraper: $(GOSRC)
	$(ENVVARS) go build $(LDFLAGS) $(BUILDFLAGS) -o $@ $(PACKAGE)
