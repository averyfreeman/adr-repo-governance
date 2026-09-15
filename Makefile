.PHONY: fmt test vet build install install-mac-layout container-binary build-container publish-container check adr-check adr-check-strict index-check detect-bs

GO ?= go
BASE ?= origin/main
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
VERSION_TAG ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
LDFLAGS ?= -X main.version=$(VERSION)
BINARY ?= ./adr
prefix ?= $(if $(PREFIX),$(PREFIX),$(HOME)/.local)
data_home ?= $(if $(XDG_DATA_HOME),$(XDG_DATA_HOME),$(prefix)/share)
bin_home ?= $(if $(XDG_BIN_HOME),$(XDG_BIN_HOME),$(prefix)/bin)
install_data_dir := $(data_home)/adr-rg
install_bindir := $(bin_home)
mac_install_data_dir := $(HOME)/Library/Application Support/org.unixgreybeard.adr-rg

DOCKER ?= docker
CONTAINER_USERNAME ?= $(or $(DOCKER_USERNAME),$(DOCKER_USER),averyfreeman)
CONTAINER_REPOSITORY ?= adr-repo-govern
CONTAINER_TAG ?= $(patsubst v%,%,$(VERSION_TAG))
CONTAINER_IMAGE ?= $(CONTAINER_USERNAME)/$(CONTAINER_REPOSITORY):$(CONTAINER_TAG)
CONTAINER_GOARCH ?= amd64
CONTAINER_PLATFORM ?= linux/$(CONTAINER_GOARCH)
CONTAINER_BINARY := ./dist/adr

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

build:
	$(GO) build -ldflags "$(LDFLAGS)" -o "$(BINARY)" ./cmd/adr

install: build
	@set -eu; \
	data_dir="$(install_data_dir)"; \
	bin_dir="$(install_bindir)"; \
	mkdir -p "$$data_dir" "$$bin_dir"; \
	test -d "$$data_dir" && test -d "$$bin_dir"; \
	cp "$(BINARY)" "$$data_dir/adr"; \
	chmod 0755 "$$data_dir/adr"; \
	if [ -e "$$bin_dir/adr" ] || [ -L "$$bin_dir/adr" ]; then \
		rm -f "$$bin_dir/adr"; \
	fi; \
	ln -s "$$data_dir/adr" "$$bin_dir/adr"

install-mac-layout: build
	@set -eu; \
	data_dir="$(mac_install_data_dir)"; \
	bin_dir="$(install_bindir)"; \
	mkdir -p "$$data_dir" "$$bin_dir"; \
	test -d "$$data_dir" && test -d "$$bin_dir"; \
	cp "$(BINARY)" "$$data_dir/adr"; \
	chmod 0755 "$$data_dir/adr"; \
	if [ -e "$$bin_dir/adr" ] || [ -L "$$bin_dir/adr" ]; then \
		rm -f "$$bin_dir/adr"; \
	fi; \
	ln -s "$$data_dir/adr" "$$bin_dir/adr"

container-binary:
	@mkdir -p "$(dir $(CONTAINER_BINARY))"
	CGO_ENABLED=0 GOOS=linux GOARCH="$(CONTAINER_GOARCH)" $(GO) build -ldflags "$(LDFLAGS)" -o "$(CONTAINER_BINARY)" ./cmd/adr

build-container: container-binary
	$(DOCKER) build --platform "$(CONTAINER_PLATFORM)" --file Dockerfile --tag "$(CONTAINER_IMAGE)" .

publish-container: build-container
	$(DOCKER) push "$(CONTAINER_IMAGE)"

adr-check: build
	$(BINARY) check adr

adr-check-strict: build
	$(BINARY) check adr --strict

index-check: build
	$(BINARY) index --check

detect-bs: build
	$(BINARY) detect-bs --base $(BASE)

check: fmt test vet build adr-check-strict index-check
