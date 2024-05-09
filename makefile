# Makefile

.PHONY: install test local container generate stepci blackbox

# Install dependencies
install:
	go get ./...

# Run Unit Testing
test:
	go test ./... -v
	
# Run the application locally
local:
	mh up mongodb
	go run src/main.go

# Build and run the Docker container
container:
	docker build --tag ghcr.io/agile-learning-institute/mentorhub-person-api:latest .
	mh up person-api
	./test/test.sh

# Generate test data
generate:
	./test/buildTestData.sh

# Run StepCI Testing
stepci:
	stepci run ./test/person.stepci.yaml

# Start containers and run stepCI testing
blackbox:
	make container
	mh up person-api
	make stepCI
