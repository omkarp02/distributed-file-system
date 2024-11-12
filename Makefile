run:
	@air

build:
	@go build -o bin/fs

start: build
	@./bin/fs

test:
	go test ./... -v