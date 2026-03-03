GO ?= /usr/local/go/bin/go
BINARY ?= ovpnctl
OUT_DIR ?= out

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X 'github.com/chalvinwz/ovpnctl/internal/cmd.Version=$(VERSION)'

.PHONY: build release clean fmt

build:
	@mkdir -p $(OUT_DIR)
	$(GO) build -o $(OUT_DIR)/$(BINARY) ./cmd/ovpnctl

release:
	@mkdir -p $(OUT_DIR)
	$(GO) build -ldflags "$(LDFLAGS)" -o $(OUT_DIR)/$(BINARY) ./cmd/ovpnctl

fmt:
	$(GO) fmt ./...

clean:
	rm -rf $(OUT_DIR)
