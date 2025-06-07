.PHONY: build install clean dev

# Build the arvia binary
build:
	go build -o bin/arvia cmd/arvia/main.go

# Install arvia globally
install: build
	go install cmd/arvia/main.go

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Development mode (for framework development)
dev:
	go run cmd/arvia/main.go serve

# Test the framework
test:
	go test ./...

help:
	@echo "Arvia Framework Build Commands:"
	@echo "  make build    - Build arvia binary"
	@echo "  make install  - Install arvia globally"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make dev      - Run in development mode"
	@echo "  make test     - Run tests"

---
