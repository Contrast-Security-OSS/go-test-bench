#!/bin/bash

function clean() {
    echo Stopping Go Demo Containers
    docker-compose -f docker-compose.demo.yml stop

    echo Removing Go Demo Containers
    docker-compose -f docker-compose.demo.yml rm
}

function rebuild() {
    echo Re-Building Go Demo Environment
    docker-compose -f docker-compose.demo.yml build --no-cache
}

if [[ "$1" == "reset" ]]
then
    clean
    rebuild
elif [[ "$1" == "update" ]] ; then
    
    clean

    echo Updating Go Test Bench
    git stash
    git pull origin

    rebuild
else
    echo Starting Go Demo Container

    docker-compose -f docker-compose.demo.yml up
fi