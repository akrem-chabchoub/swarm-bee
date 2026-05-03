.PHONY: build run test clean

build:
	go build -o bin/swarm-bee ./cmd/swarm-bee

run:
	go run ./cmd/swarm-bee

test:
	go test -v -race -count=1 ./...

clean:
	rm -rf bin/