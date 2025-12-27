.PHONY: build test clean

build:
	go build -o kex ./cmd/kex

test:
	go test -v ./...

clean:
	rm -f kex
