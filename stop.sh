echo Stopping Go Demo Container

docker-compose -f docker-compose.demo.yml stop
docker-compose -f docker-compose.demo.yml rm
docker-compose -f docker-compose.demo.yml build --no-cache

