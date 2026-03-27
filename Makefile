generate:
	go generate ./...

test: unit-test lint

lint:
	golangci-lint run -v

unit-test:
	go tool ginkgo -r .
