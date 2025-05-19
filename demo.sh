#!/bin/bash

if [[ "$1" == "gin" ]]; then
    echo "Starting Go Demo Container for Gin framework"
elif [[ "$1" == "julienschmidt" ]]; then
    echo "Starting Go Demo Container for Julienschmidt/httprouter framework"
elif [[ "$1" == "go-swagger" ]]; then
    echo "Starting Go Demo Container for Swagger framework"
elif [[ "$1" == "chi" ]]; then
    echo "Starting Go Demo Container for Chi v5 framework"
else
    echo "Starting Go Demo Container for standard library"
    docker build -f Dockerfile.agent -t contrast-go-demo .
    docker run -p 8080:8080 contrast-go-demo
    exit 0
fi
docker build --build-arg FRAMEWORK=$1 -f Dockerfile.agent -t contrast-go-demo .
docker run -p 8080:8080 contrast-go-demo
