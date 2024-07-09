
GOLANGCI_LINT_VERSION := v1.59.1

lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin $(GOLANGCI_LINT_VERSION)
	bin/golangci-lint run --timeout 2m

test:
	go test --cover ./...

