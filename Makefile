
.PHONY: deps
deps:
	@echo "installing Mockery ..."
	@brew install mockery

.PHONY: mocks
mocks:
	@mockery --dir ./internal/storage --name Database --output ./mocks/database --outpkg database --with-expecter