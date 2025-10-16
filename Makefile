.POSIX:

bin/:
	mkdir -p $@

bin/sqlc: | bin/
	GOBIN="$(realpath $(dir $@))" go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: sqlc
sqlc: | bin/sqlc
	./bin/sqlc generate -f internal/db/pg/sqlc.yaml

bin/golangci-lint: | bin/
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh

.PHONY: lint
lint: bin/golangci-lint
	go list -f '{{.Dir}}/...' -m | xargs ./bin/golangci-lint run --fix

.PHONY: local
local:
	docker-compose -f .docker/local/docker-compose.yml up --build -d
