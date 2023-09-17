set dotenv-load := true

start:
    go run example/server.go

test:
    go test `go list ./... | grep -v mocks`

test-cov:
    go test `go list ./... | grep -v mocks` -covermode=count -coverprofile cover.out
    go tool cover -func cover.out

generate:
    go generate ./...

lint:
	golangci-lint run