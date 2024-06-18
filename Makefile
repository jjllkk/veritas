GOARCH=amd64
GO111MODULE=on

.PHONY: build

build: wasm

wasm:
	GOOS=js GOARCH=wasm go build -o web/veritas.wasm ./cmd/wasm

run:
	go run -tags server ./cmd/server

all: build run

test:
	go test ./...

clean:
	rm -rf web/veritas.wasm

