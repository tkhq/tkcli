.PHONY: all
all: build

.PHONY: build
build:
	go build -o build/tk cmd/tk/main.go

.PHONY: clean
clean:
	rm -rf build/

.PHONY: test
test:
	go test ./...
