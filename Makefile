.PHONY: help
help: ## show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: deps
deps: ## install Mockery
	@echo "installing Mockery ..."
	@brew install mockery

.PHONY: fmt
fmt: ## format *.go files
	go fmt ./...

.PHONY: test
test: ## run all tests
	go test ./... -cover

.PHONY: vendor
vendor: ## save all deps to vendor
	go mod vendor

.PHONY: mocks
mocks: ## rebuild all mocks
	@mockery --dir ./internal/storage --name Database --output ./mocks/database --outpkg database --with-expecter
	@mockery --dir ./internal/importer --name DataSaver --output ./mocks/importer --outpkg importer --with-expecter
	@mockery --dir ./internal/gateway --name DAO --output ./mocks/gateway --outpkg gateway --with-expecter