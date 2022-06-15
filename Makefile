.PHONY: test
test:
	go test ./... -race -v

.PHONY: install-mockgen
install-mockgen:
	go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: generate
generate:
	go generate ./...