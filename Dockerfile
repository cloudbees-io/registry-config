FROM golang:1.22-alpine3.20 AS build
COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download
COPY . /src/
ENV CGO_ENABLED=0
ARG VERSION=dev
RUN go build -a -ldflags '-X main.version=$VERSION -w -extldflags "-static"' -o registry-config ./cmd/registry-config

FROM gcr.io/distroless/static:nonroot
COPY --from=build /src/registry-config /registry-config
COPY ./pkg/registries/testdata/registries.json /etc/cloudbees/registries.json
ENTRYPOINT ["/registry-config"]
