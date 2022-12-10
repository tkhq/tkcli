.PHONY: all
all: build

.PHONY: local-release
local-release:
	goreleaser release --snapshot --rm-dist

.PHONY: test
test: build/turnkey
	go test ./...

build/turnkey:
	go build -o build/turnkey main.go

.PHONY: clean
clean:
	rm -rf dist/ build/