.PHONY: test
test: mod-tidy
	go test -v ./...

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build
build: test
	go build ./cmd/ccv
