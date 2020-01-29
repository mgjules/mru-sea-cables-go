.PHONY: build

build:
	@go build -mod="vendor" -ldflags="-s -w"