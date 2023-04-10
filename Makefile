.PHONY: all
all: build/turnkey

.PHONY: local-release
local-release:
	goreleaser release --snapshot --rm-dist

.PHONY: test
test: build/turnkey
	go test ./...

build/turnkey: main.go internal/
	go build -o build/turnkey main.go

.PHONY: clean
clean:
	rm -rf dist/ build/
