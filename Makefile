PKG=github.com/redish101/depositum
SERVER=./cmd/depositum
PNPM=pnpm --dir ui

all: server

server: ui
	go build -v -o bin/depositum $(SERVER)

apidoc:
	swag init -g cmd/depositum/main.go -d . --parseInternal --parseDependency
	swag fmt

clean:
	rm -r bin

ui:
	$(PNPM) build
	rm -rf internal/ui/client
	cp -r ui/build/client internal/ui/client

fmt:
	go fmt ./...
	goimports -w .
	$(PNPM) format

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: all clean fmt server test coverage ui