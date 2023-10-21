cd "$(find ~ -name "institute-person-api" | head -n 1)"
./docker-build.sh

cd "$(find ~ -name "institute-person-api" | head -n 1)"
cd /src/docker
./docker-build.sh

docker image prune -f

docker-compose up --detach