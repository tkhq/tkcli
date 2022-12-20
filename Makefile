.PHONY: all
all: build

.PHONY: local-release
local-release:
	goreleaser release --snapshot --rm-dist

.PHONY: test
test: build
	go test ./...

.PHONY: build
build:
	go build -o build/turnkey main.go

.PHONY: clean
clean:
	rm -rf dist/ build/
