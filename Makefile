.PHONY: build clean test lint install

BINARY=kopy

build:
	go build -o bin/$(BINARY) .

clean:
	rm -rf bin/ dist/

test:
	go test ./... -v -cover

lint:
	golangci-lint run ./...

install: build
	cp bin/$(BINARY) $(GOPATH)/bin/$(BINARY)

run: build
	./bin/$(BINARY) --help

cross:
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY)-windows-amd64.exe .

tidy:
	go mod tidy
