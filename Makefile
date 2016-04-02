all: build test install

deps:
	godep save ./...

build: deps
	godep go build ./...

test: deps
	godep go test . ./cmd/...

install: deps
	godep go install ./...