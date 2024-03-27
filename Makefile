#!/usr/bin/make -f
export GOPROXY=https://goproxy.cn,direct
build: go.sum
ifeq ($(OS),Windows_NT)
	go build  -o build/farm.exe ./cmd/farm
else
	go build  -o build/farm ./cmd/farm
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

install:
	go install  ./cmd/farm

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs goimports -w -local github.com/irisnet/coinswap-server