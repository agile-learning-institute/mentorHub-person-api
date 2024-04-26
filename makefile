# Makefile

.PHONY: test generate

# Run the application locally
run:
	mh up mongodb
	go run src/main.go

# Build and run the Docker container
build:
	docker build --tag ghcr.io/agile-learning-institute/mentorhub-person-api:latest .
	mh up person-api

# Run tests
test:
	./test/test.sh

# Generate test data
generate:
	./test/buildTestData.sh
