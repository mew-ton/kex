.PHONY: build test clean init e2e

build:
	go build -o bin/kex ./cmd/kex

test:
	go test -v ./...

e2e:
	go test -v ./e2e/...

clean:
	rm -f bin/kex

init:
	go tool lefthook install
