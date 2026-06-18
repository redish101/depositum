PKG=github.com/redish101/depositum
SERVER=./cmd/depositum

all: server

server:
	go build -v $(SERVER)

clean:
	rm server

fmt:
	go fmt ./...
	goimports -w .

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: all clean fmt server test coverage