GO ?= /usr/local/go/bin/go
BINARY ?= ovpnctl
OUT_DIR ?= out

.PHONY: build test cover clean fmt

build:
	@mkdir -p $(OUT_DIR)
	$(GO) build -o $(OUT_DIR)/$(BINARY) ./cmd/ovpnctl

test:
	$(GO) test ./...

cover:
	$(GO) test -coverprofile=cover.out ./...
	$(GO) tool cover -func=cover.out

fmt:
	$(GO) fmt ./...

clean:
	rm -rf $(OUT_DIR) cover.out
