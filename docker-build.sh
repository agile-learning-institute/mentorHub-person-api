GOOS=linux GOARCH=amd64 go build -o "institute-person-api" main.go
export BRANCH=$(git branch --show-current)
export PATCH=$(git rev-parse $BRANCH)
echo $BRANCH.$PATCH > PATCH_LEVEL
docker build . --tag ghcr.io/agile-learning-institute/institute-person-api:latest
# docker tag institute-person-api:latest ghcr.io/agile-learning-institute/institute-person-api:latest
# docker push ghcr.io/agile-learning-institute/institute-person-api:latest
