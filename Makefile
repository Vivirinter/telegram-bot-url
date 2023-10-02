BINARY_NAME=telegram-bot-url
SRC_DIR=cmd/main
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go')

GO=go
GO_BUILD=$(GO) build
GO_CLEAN=$(GO) clean
GO_TEST=$(GO) test

all: build
build: test
		$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v ./cmd/main

test:
		$(GO_TEST) -v ./...

clean:
		$(GO_CLEAN)
		rm -f $(BUILD_DIR)/*