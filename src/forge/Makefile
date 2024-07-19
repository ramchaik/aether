# Directories
PROTO_DIR := ../../protobuf
GO_OUT_DIR := ./internal/genprotobuf

# Find all .proto files
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Extract package names from .proto files
PACKAGES := $(sort $(shell grep "^package" $(PROTO_FILES) | awk '{print $$2}' | sed 's/;//'))

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/worker/main.go

# Run the application
run:
	@go run cmd/worker/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary and generated protobuf files
clean:
	@echo "Cleaning..."
	@rm -f main
	@rm -rf $(GO_OUT_DIR)

# Live Reload
watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/cosmtrek/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

generate-protobuf: $(PACKAGES)

$(PACKAGES):
	@echo "Generating protobuf files for package: $@"
	@mkdir -p $(GO_OUT_DIR)/$@
	@protoc -I=$(PROTO_DIR) \
		--go_out=$(GO_OUT_DIR)/$@ --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR)/$@ --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/$@.proto

list-packages:
	@echo "Available packages:"
	@echo $(PACKAGES) | tr ' ' '\n'

.PHONY: all build run test clean generate-protobuf list-packages $(PACKAGES)