export GO111MODULE=on
export GOPROXY=https://proxy.golang.org
export GOSUMDB=off

LOCAL_BIN:=$(CURDIR)/bin
BUILD_ENVPARMS:=CGO_ENABLED=0

GOLANGCI_TAG:=1.48.0

export PATH:=$(LOCAL_BIN):$(PATH)

.PHONY: deps
deps:
	$(info Installing binary dependencies...)
	go mod download
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCI_TAG)
	GOBIN=$(LOCAL_BIN) go install mvdan.cc/gofumpt@latest
	GOBIN=$(LOCAL_BIN) go install github.com/incu6us/goimports-reviser/v3@latest
	npm install cspell --prefix $(LOCAL_BIN)
	go mod tidy

.PHONY: clean
clean:
	rm -rf bin

.PHONY: build
build:
	$(info Building...)
	$(BUILD_ENVPARMS) go build -o ./bin/port-service ./cmd/port-api

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint: 
	$(info Running lint...)
	GOBIN=$(LOCAL_BIN) golangci-lint run  --config=.cfg/lint.yaml ./...

.PHONY: format
format:
	GOBIN=$(LOCAL_BIN) goimports-reviser -project-name port-service -company-prefixes github.com/canbo-x/port-service ./...
	GOBIN=$(LOCAL_BIN) gofumpt -l -w -extra .

.PHONY: spell
spell:
	$(info Spell checking for the source code...)
	$(LOCAL_BIN)/node_modules/cspell/bin.js "**" --config ./cspell.yml 