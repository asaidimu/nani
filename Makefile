.PHONY: all build test clean

all: build

build:
	@mkdir -p ./dist
	go build -o ./dist/main ./main.go

test:
	go test -v ./...

clean:
	@rm -f ./dist/main
	@rmdir -p ./dist 2>/dev/null || true

