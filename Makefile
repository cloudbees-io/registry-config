
VERSION ?= dev

GOLANGCI_LINT_VERSION := v1.59.1


all: cli

cli:
	CGO_ENABLED=0 go build -a -ldflags '-X main.version=$(VERSION) -w -extldflags "-static"' -o ./bin/registry-config ./cmd/registry-config

container:
	docker build -t registry-config --build-arg VERSION=$(VERSION) .

lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin $(GOLANGCI_LINT_VERSION)
	bin/golangci-lint run --timeout 2m

test:
	go test --cover ./...
