# registry-config

Container registry configuration management for use in pipelines on the CloudBees platform

## Usage examples

The CLI built within this repo is published as a container image that you could use as is as a step within an Action or copy the contained binary into your own image.
The following examples assume you have this git repo checked out and the binary built locally.

### Resolving an image

Resolve an image reference:
```
./bin/registry-config --config=./pkg/registries/testdata/registries.json resolve golang:1.22
```

Alternatively, write the resolved image references into a file:
```
./bin/registry-config --config=./pkg/registries/testdata/registries.json resolve golang:1.22 outfile.txt
```

### Converting the configuration

To convert the CloudBees registries configuration JSON file into [Red Hat's registries.conf format](https://github.com/containers/image/blob/v5.31.0/docs/containers-registries.conf.5.md#example), run:
```
./bin/registry-config --config=./pkg/registries/testdata/registries.json convert ./registries.conf
```

## Development

Build the binary:
```
make cli
```

Run tests:
```
make test
```

Run the linter:
```
make lint
```

Build the container image:
```
make container
```
