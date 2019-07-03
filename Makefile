generate:
	go generate ./...

test: unit-test lint

lint:
	golangci-lint run -v

unit-test:
	ginkgo -r .
