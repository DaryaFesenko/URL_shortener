APP = url_shortener

HAS_LINT := $(shell command -v golangci-lint;)
HAS_IMPORTS := $(shell command -v goimports;)

.PHONY: all
all: run

.PHONY: lint
lint: bootstrap
	@echo "+ $@"
	@golangci-lint run

.PHONY: run
run: clean build
	@echo "+ $@"
	./${APP} -c "./cmd/url_shortener/config.yaml"

.PHONY: build
build: lint
	@echo "+ $@"
	@go build ./cmd/url_shortener/

.PHONY: clean
clean:
	@echo "+ $@"
	@rm -f ./${APP}

.PHONY: bootstrap
bootstrap:
	@echo "+ $@"
ifndef HAS_LINT
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
endif
ifndef HAS_IMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif

.PHONY: test
test: 
	@go test -v -coverprofile cover.out ./...
	