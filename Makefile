.PHONY: build

buildrun: lint build run

lint:
	@golangci-lint run

# Do NOT run too often (i.e rate-limiting GitHub)
license: build
	@golicense ./golicense.hcl mru-sea-cables-go

run:
	@./mru-sea-cables-go

build:
	@go build -mod="vendor" -ldflags="-s -w"