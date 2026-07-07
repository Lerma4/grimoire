# Grimoire Makefile — common dev tasks.
BINARY  := bin/grimoire
PKG     := ./cmd/grimoire
GOFLAGS :=

.PHONY: all run build install test fmt vet lint clean check

all: build

## run: launch the TUI via go run
run:
	go run $(GOFLAGS) $(PKG)

## build: compile the binary into ./bin
build:
	go build $(GOFLAGS) -o $(BINARY) $(PKG)

## install: install the global `grimoire` command via go install
install:
	go install $(GOFLAGS) $(PKG)

## test: run the test suite
test:
	go test ./...

## fmt: format all Go sources
fmt:
	gofmt -w .

## vet: run go vet
vet:
	go vet ./...

## check: fmt + vet + test (the pre-commit gate)
check: fmt vet test

## lint: run golangci-lint if installed (non-fatal if missing)
lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || \
		echo "golangci-lint not installed; skipping (see README to install)"

## clean: remove build artifacts
clean:
	rm -rf bin/ coverage.txt
