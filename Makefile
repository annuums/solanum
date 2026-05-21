
.PHONY: test
test:
	go test -v ./... -short

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...
