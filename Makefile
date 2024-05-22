.PHONY: test
test: mod-tidy
	export GIT_CONFIG_NOSYSTEM=true && \
		go test -v ./...

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build
build:
	GOVERSION=$$(go version) \
						goreleaser build --clean --debug --single-target --snapshot

.PHONY: lint
lint:
	golangci-lint run --enable gocritic

.PHONY: cover
cover: mod-tidy generate
	go test -v -covermode=atomic -coverprofile=cover.out.raw -coverpkg=./... ./...
	grep -Ev 'internal/mock|_enumer.go' cover.out.raw > cover.out
	go tool cover -html=cover.out
