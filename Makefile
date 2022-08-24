
.PHONY: deps
deps:
	@echo "installing Mockery ..."
	@brew install mockery

.PHONY: mocks
mocks:
	@mockery --dir ./internal/storage --name Database --output ./mocks/database --outpkg database --with-expecter
	@mockery --dir ./internal/importer --name DataSaver --output ./mocks/importer --outpkg importer --with-expecter
	@mockery --dir ./internal/gateway --name DAO --output ./mocks/gateway --outpkg gateway --with-expecter