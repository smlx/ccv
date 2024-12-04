.PHONY: test
test: mod-tidy
	export GIT_CONFIG_NOSYSTEM=true GIT_CONFIG_GLOBAL=/tmp/gitconfig \
		&& git config --global user.email "test@example.com" \
		&& git config --global user.name "Test" \
		&& go test -v ./... -count=1

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build
build:
	goreleaser build --clean --debug --single-target --snapshot

.PHONY: lint
lint:
	golangci-lint run --enable gocritic

.PHONY: cover
cover: mod-tidy generate
	go test -v -covermode=atomic -coverprofile=cover.out.raw -coverpkg=./... ./...
	grep -Ev 'internal/mock|_enumer.go' cover.out.raw > cover.out
	go tool cover -html=cover.out
