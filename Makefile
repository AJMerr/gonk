GO ?= go
PKG := ./...
BINARY := ./bin/gonk
MAIN := ./cmd/gonk


build:
	$(GO) build -o $(BINARY) -ldflags "-X main.Version=1.0.0" $(MAIN)

test:
	$(GO) test $(PKG) -race $(ARGS)

tidy:
	$(GO) mod tidy

fmt: 
	$(GO) fmt $(PKG)

vet: 
	$(GO) vet $(PKG)