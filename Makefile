GOCMD?=go
GOFLAGS?=
GOTESTFLAGS?=-v -count=1 -parallel=1
YANEURAGO_BIN?=yaneurago

build: out
	@echo "> Building..."
	$(GOCMD) build -o ./out/$(YANEURAGO_BIN) -a ./cmd/yaneurago

out:
	@mkdir out

install-golangci-lint:
	@echo "> Installing golangci-lint..."
	cd tools && $(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-golangci-lint
	@echo "> Linting code..."
	@golangci-lint run -c golangci.yaml

.PHONY: lint build