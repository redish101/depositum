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

.PHONY: all clean fmt depositum test