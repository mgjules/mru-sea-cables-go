.PHONY: build

buildrun: build run

run:
	@./mru-sea-cables-go

build:
	@go build -mod="vendor" -ldflags="-s -w"