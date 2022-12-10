.PHONY: all
all: build

.PHONY: build
build:
	goreleaser release --snapshot --rm-dist

.PHONY: clean
clean:
	rm -rf dist/

.PHONY: test
test:
	go test ./...
