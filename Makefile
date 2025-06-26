.PHONY: test

test:
	@echo "Running tests..."
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out