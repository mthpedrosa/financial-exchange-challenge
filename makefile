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

.PHONY: mocks
mocks:
	rm -fr ./mocks
	mockery --all --dir ./internal --disable-version-string --case snake --keeptree
	
.PHONY: swag
swag:
	rm -fr ./docs
	swag init -g cmd/main.go -parseDependency docs

# Makefile
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

migrate-up:
	migrate -database "$(DATABASE_URL)" -path db/migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path db/migrations down