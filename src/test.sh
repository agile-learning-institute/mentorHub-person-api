#!/bin/bash

# Default port
port=${1:-8082}

# Function to test HTTP request and exit if failed
test_http_request() {
    local url="$1"
    local description="$2"
    local method="$3"
    local payload="$4"

    echo "Testing $description"
    local response_code
    if [ -z "$payload" ]; then
        response_code=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "localhost:$port$url")
    else
        response_code=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" -d "$payload" "localhost:$port$url")
    fi

    if [ "$response_code" -eq 200 ]; then
        echo "$description successful"
    else
        echo "$description failed. Received HTTP code $response_code"
        exit 1
    fi
}

# Main function
main() {
    local current_time
    current_time=$(date +"%Y-%m-%d_%H-%M-%S")

    test_http_request "/api/config/" "Get Config" "GET"
    test_http_request "/api/health/" "Get health" "GET"
    test_http_request "/api/enums/" "Get enums" "GET"
    test_http_request "/api/partners/" "Get partners" "GET"
    test_http_request "/api/mentors/" "Get mentors" "GET"
    test_http_request "/api/person/" "Get people" "GET"
    test_http_request "/api/person/aaaa00000000000000000017" "Get person" "GET"
    test_http_request "/api/person/aaaa00000000000000000017" "Update Person" "PATCH" '{"name":"Foo", "description":"Some short description"}'
    test_http_request "/api/person/" "Create Person" "POST" "{\"name\":\"Person_$current_time\", \"description\":\"A New Person\"}"

    echo "SUCCESS!!!!"
}

# Execute main function
main "$@"
