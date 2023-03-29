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

.PHONY: api
api:
	swagger generate client -f https://raw.githubusercontent.com/tkhq/sdk/main/packages/http/src/__generated__/services/coordinator/public/v1/public_api.swagger.json -t api
	go mod tidy

.PHONY: lint
lint:
	go vet ./...

.PHONT: install
install: all
	go install ./cmd/turnkey

.PHONY: build/turnkey
build/turnkey: cmd/turnkey/ internal/
	go build -o build/turnkey ./cmd/turnkey

.PHONY: clean
clean:
	rm -rf dist/ build/
