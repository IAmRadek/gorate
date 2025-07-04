BINARY_NAME=gorate

# Build directory
BUILD_DIR=build

# Main application directory
MAIN_DIR=cmd/gorate

# Environment file
ENV_FILE=.development.env

# Check if env file exists
ifneq ("$(wildcard $(ENV_FILE))","")
	include $(ENV_FILE)
	export
endif

.PHONY: all build clean run env-check

# Add env-check to ensure .development.env exists
env-check:
	@if [ ! -f "$(ENV_FILE)" ]; then \
		echo "Error: $(ENV_FILE) does not exist"; \
		exit 1; \
	fi

build-binary-in-docker:
	@echo "Building..."
	@CGO_ENABLED=0 GOOS=linux go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_DIR)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)


run: env-check
	@echo "Running with $(ENV_FILE)..."
	@go run -race ./$(MAIN_DIR)

tests:
	@echo "Running tests"
	@go test ./... -v

env: env-check
	@echo "Environment variables from $(ENV_FILE):"
	@cat $(ENV_FILE)

build-dockerimage:
	docker build -t gorate .

run-docker:
	docker run --rm --name gorate \
  		--env-file $(ENV_FILE) \
  		-p 8080:8080 \
  		gorate:latest
