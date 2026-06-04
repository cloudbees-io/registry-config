FROM gcr.io/distroless/static:nonroot
COPY registry-config /registry-config
COPY ./pkg/registries/testdata/registries.json /etc/cloudbees/registries.json
ENTRYPOINT ["/registry-config"]
