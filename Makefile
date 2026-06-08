all: depositum

depositum:
	go build -v

clean:
	rm depositum

fmt:
	go fmt ./...
	goimports -w .

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: all clean fmt depositum test coverage