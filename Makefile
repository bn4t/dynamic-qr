SHELL = /bin/bash

# use bash strict mode
.SHELLFLAGS := -eu -o pipefail -c

.ONESHELL:
.DELETE_ON_ERROR:

.SUFFIXES:      # delete the default suffixe
.SUFFIXES: .go  # add .go as suffix

PREFIX?=/usr/local
_INSTDIR=$(DESTDIR)$(PREFIX)
BINDIR?=$(_INSTDIR)/bin
GO?=go
GOFLAGS?=
STATIK?=~/go/bin/statik
RM?=rm -f # Exists in GNUMake but not in NetBSD make and others.


all: build

build:
	$(STATIK) -f -src="./static/" -dest="internal/" -p="statik"
	$(GO) build $(GOFLAGS) -o dynamic-qr .

run: build
	./dynamic-qr

clean:
	$(RM) dynamic-qr


.DEFAULT_GOAL = all
.PHONY: all build install uninstall clean release

