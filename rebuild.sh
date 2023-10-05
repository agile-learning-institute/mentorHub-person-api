GOOS=linux GOARCH=amd64 go build -o "institute-person-api" main.go
docker build . --tag institute-person-api
docker compose up --detach      
clear; docker logs -f institute-person-api-institute-person-api-1