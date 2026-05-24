all: depositum

depositum:
	go build -v

clean:
	rm depositum

fmt:
	goimports -w .

PHONY: all clean fmt