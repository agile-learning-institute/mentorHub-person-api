#!/bin/bash

# Default port
sequence=${1:-1}
testAll=${2:-0}
port=8082

# Function to test HTTP request and exit if failed
test_http_request() {
    local url="$1"
    local method="$2"
    local payload="$3"

    local response_code
    if [ -z "$payload" ]; then
        response_code=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "localhost:$port$url")
    else
        response_code=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" -d "$payload" "localhost:$port$url")
    fi

    if [ "$response_code" -ne 200 ]; then
        echo "$description failed. Received HTTP code $response_code"
        exit 1
    fi
}

# Main function
main() {
    local current_time
    firstName=$(sort -R ./test/loadData/firstNames.txt | head -n 1 | tr -d '\n')
    lastName=$(sort -R ./test/loadData/lastNames.txt | head -n 1 | tr -d '\n')
    status=$(sort -R ./test/loadData/status.txt | head -n 1 | tr -d '\n')
    cadence=$(sort -R ./test/loadData/cadence.txt | head -n 1 | tr -d '\n')
    device=$(sort -R ./test/loadData/device.txt | head -n 1 | tr -d '\n' )
    title=$(sort -R ./test/loadData/title.txt | head -n 1 | tr -d '\n' )
    roles=$(sort -R ./test/loadData/roles.txt | head -n 1 | tr -d '\n' )
    notes=$(sort -R ./test/loadData/BOFH.txt | head -n 1 | tr -d '\n\t\r"' | tr -d "'" | cut -c 1-255)
    mentor=$(sort -R ./test/loadData/mentors.txt | head -n 1 | tr -d '\n\t\r')
    partner=$(sort -R ./test/loadData/partners.txt | head -n 1 | tr -d '\n\t\r')
    echo $sequence, $firstName, $lastName, $status, $cadence, $device, $title, $roles, $notes

    if [ "$testAll" -eq 0 ]; then
        test_http_request "/api/config/" "GET"
        test_http_request "/api/health/" "GET"
        test_http_request "/api/person/" "GET"
        test_http_request "/api/person/aaaa00000000000000000017" "GET"
        test_http_request "/api/person/aaaa00000000000000000017" "PATCH" '{"userName":"Foo", "description":"Some short description"}'
    fi

    test_http_request "/api/person/" \
        "POST" \
        "{\"userName\":\"$firstName.$lastName.$sequence\", \
        \"firstName\":\"$firstName\", \
        \"lastName\":\"$lastName\", \
        \"eMail\":\"$firstName.$lastName@gmail.com\", \
        \"gitHub\":\"$lastName.$lastName\", \
        \"phone\":\"888 555-1212\", \
        \"status\":\"$status\", \
        \"roles\":$roles, \
        \"device\":\"$device\", \
        \"title\":\"$title\", \
        \"description\":\"$notes\", \
        \"location\":\"Somewhere Over Yonder\", \
        \"mentorId\":\"$mentor\", \
        \"partnerId\":\"$partner\", \
        \"cadence\":\"$cadence\"}"
}

# Execute main function
main "$@"
