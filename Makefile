.PHONY: fmt test build check adr-check adr-check-strict bs-detector

BASE ?= origin/main

fmt:
	gofmt -w .

test:
	go test ./...

build:
	go build ./cmd/adr

adr-check:
	./adr check adr

adr-check-strict:
	./adr check adr --strict

bs-detector:
	./adr bs-detector --base $(BASE)

check: fmt test build adr-check-strict
