.PHONY: all
all: fmt lint build

.PHONY: local-release
local-release:
	goreleaser release --snapshot --rm-dist

.PHONY: test
test: build/turnkey
	go test ./...

build: build/turnkey

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: build/turnkey
build/turnkey: main.go internal/
	go build -o build/turnkey .

.PHONY: clean
clean:
	rm -rf dist/ build/
