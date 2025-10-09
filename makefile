.PHONY: lint
lint:
	golangci-lint -v run

.PHONY: format
format:
	goimports -w ./
	gofmt -s -w ./
	swag fmt

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: cover
cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./internal/...
	gocov coverage.txt
	rm coverage.txt