wire:
	@wire gen ./di/container

gql:
	@go run github.com/99designs/gqlgen generate

lint:
	@golangci-lint run

start:
	@ZITADEL_KEY_PATH=assets/local/zitadel.json go run main.go start

freshstart:
	@make wire && make gql && make start

test:
	@go test -failfast -v -covermode=set -coverpkg=./internal/... -coverprofile cover.out ./tests/integration/...

testcoverage:
	@make test && go-test-coverage --config=./.testcoverage.yml

testcoveragereport:
	@make test && go tool cover -html cover.out -o cover.html