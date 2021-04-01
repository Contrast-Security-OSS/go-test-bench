# Contrast Agent Demonstration

## Step 1 - Prerequisites

Install Git and Docker to your system

* [Git](https://git-scm.com/)
* [Docker](https://docs.docker.com/get-docker/)

## Step 2 - Clone the go-test-bench application

Open the terminal on your local system and navigate to a safe working directory.

`cd ~`

Clone the Contrast Security `go-test-bench` to your local system using Git.

`git clone https://github.com/Contrast-Security-OSS/go-test-bench.git`

Move into the go-test-bench directory.

`cd go-test-bench`

## Step 3 - Copy contrast_security.yaml to ~/go-test-bench

`cp /path/to/config/contrast_security.yaml .`

## Step 4 - Build the Docker Compose images

`docker-compose -f docker-compose.demo.yml build --no-cache`

## Step 5 - Start the instrumented `go-test-bench`

`docker-compose -f docker-compose.demo.yml up`

## Step 6 - Navigate to the go-test-bench UI

[Click Here to visit the Test Bench - http://localhost:8080](http://localhost:8080)

## Cleanup

To ensure a clean and updated environment for demonstrations execute the following commands.

```bash
docker-compose -f docker-compose.demo.yml stop
docker-compose -f docker-compose.demo.yml rm
docker-compose -f docker-compose.demo.yml build --no-cache
```

This will re-build your test bench with the most recent copies of the go-agent, contrast-service and test bench.
