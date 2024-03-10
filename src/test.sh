#!/bin/bash

# Default port
sequence=${1:-1}
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
    firstName=$(sort -R ./loadData/firstNames.txt | head -n 1 | tr -d '\n')
    lastName=$(sort -R ./loadData/lastNames.txt | head -n 1 | tr -d '\n')
    status=$(sort -R ./loadData/status.txt | head -n 1 | tr -d '\n')
    cadence=$(sort -R ./loadData/cadence.txt | head -n 1 | tr -d '\n')
    device=$(sort -R ./loadData/device.txt | head -n 1 | tr -d '\n' )
    title=$(sort -R ./loadData/title.txt | head -n 1 | tr -d '\n' )
    roles=$(sort -R ./loadData/roles.txt | head -n 1 | tr -d '\n' )
    notes=$(sort -R ./loadData/BOFH.txt | head -n 1 | tr -d '\n\t\r"' | tr -d "'" | cut -c 1-255)
    mentor=$(sort -R ./loadData/mentors.txt | head -n 1 | tr -d '\n\t\r')
    partner=$(sort -R ./loadData/partners.txt | head -n 1 | tr -d '\n\t\r')
    echo $sequence, $firstName, $lastName, $status, $cadence, $device, $title, $roles, $notes

    test_http_request "/api/config/" "GET"
    test_http_request "/api/health/" "GET"
    test_http_request "/api/person/" "GET"
    test_http_request "/api/person/aaaa00000000000000000017" "GET"
    test_http_request "/api/person/aaaa00000000000000000017" "PATCH" '{"name":"Foo", "description":"Some short description"}'
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
