cd ../institute-mongodb
docker build . --tag institute-mongosh

cd ../institute-person-api
GOOS=linux GOARCH=amd64 go build -o "institute-person-api" main.go
export BRANCH=$(git branch --show-current)
export PATCH=$(git rev-parse $BRANCH)
echo $BRANCH.$PATCH > PATCH_LEVEL
docker build . --tag institute-person-api

docker image prune -f

docker compose up --detach