.DEFAULT_GOAL := help

PROJECT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
OUTPUT_PATH := $(PROJECT_PATH)/_output

TAG ?=  $(shell git -C "$(PROJECT_PATH)" rev-parse HEAD)
ifdef VERSION
override TAG = $(VERSION)
endif

GO_VERSION=1.24
GO_IMPORTS_VERSION=latest
GOLANG_CI_LINT_VERSION=v2.2.2
MOCKERY_VERSION=3.5

.PHONY: unit
unit: ## Run unit tests
	mkdir -p "$(OUTPUT_PATH)"
	go install github.com/jstemmer/go-junit-report/v2@v2.1.0
	GOEXPERIMENT=nocoverageredesign go test ./... -v -covermode=count -coverprofile=$(OUTPUT_PATH)/cover.out | tee /dev/stderr | go-junit-report -set-exit-code > $(OUTPUT_PATH)/report.xml

.PHONY: unit-local
unit-local: ## Run unit tests locally
	mkdir -p "$(OUTPUT_PATH)"
	GOEXPERIMENT=nocoverageredesign go test ./... -v -covermode=count -coverprofile=$(OUTPUT_PATH)/cover.out | tee /dev/stderr | go-junit-report -set-exit-code > $(OUTPUT_PATH)/report.xml

.PHONY: unit-coverage
unit-coverage: unit ## Runs unit tests and generates a html coverage report
	go tool cover -html=$(OUTPUT_PATH)/cover.out

.PHONY: unit-coverage-cobertura
unit-coverage-cobertura: ## Runs unit tests and generates a cobertura coverage report
	go install github.com/boumenot/gocover-cobertura@v1.2.0
	gocover-cobertura < $(OUTPUT_PATH)/cover.out > $(OUTPUT_PATH)/coverage.xml

fix-lint: ## Automatically fixes formatting issues (goimports, gofmt)
	docker run --rm \
		-v $(PROJECT_PATH):/app \
		-w /app \
		golang:$(GO_VERSION) \
		bash -c "\
			go install golang.org/x/tools/cmd/goimports@latest && \
			goimports -w -format-only . && \
			gofmt -s -w ."

.PHONY: lint
lint: ## Checks code formatting/quality
	docker run --rm -v $(PROJECT_PATH):/app -v ~/.gitconfig:/etc/gitconfig -w /app golangci/golangci-lint:$(GOLANG_CI_LINT_VERSION) golangci-lint run -v --timeout=5m

help:
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
