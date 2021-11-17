# Contrast Agent Demonstration

## Step 1 - Prerequisites

Install Git and Docker to your system

* [Git](https://git-scm.com/)
* [Docker](https://docs.docker.com/get-docker/)

Make sure the Docker daemon is running. You can check this by running

`docker info`

If that command succeeds, then you're ready to continue to the next step.

## Step 2 - Clone the go-test-bench application

Open the terminal on your local system and navigate to a safe working directory.

`cd ~`

Clone the Contrast Security `go-test-bench` to your local system using Git.

`git clone https://github.com/Contrast-Security-OSS/go-test-bench.git`

Move into the go-test-bench directory.

`cd go-test-bench`

## Step 3 - Copy contrast_security.yaml to the go-test-bench directory

Download your configuration YAML file from the Contrast environment by selecting the "Add agent" button.

`cp /path/to/config/contrast_security.yaml .`

## Step 4 - Start the DEMO environment

`./demo.sh $FRAMEWORK`

example: `./demo.sh std`

The Go Test Bench supports the following frameworks:

* std (standard library `net/http`)
* [gin](https://github.com/gin-gonic/gin)

If no framework is specified, the standard library is used.

**NOTE: The first execution will take time because it has to build the environment**

## Step 5 - Navigate to the go-test-bench UI

[Click Here to visit the Test Bench - http://localhost:8080](http://localhost:8080)

---

## Maintaining your DEMO Environment

### Updating

To ensure an updated environment for demonstrations execute the following command from the go-test-bench directory.

```bash
./demo.sh update
```

### Resetting your DEMO environment

To ensure a clean and updated environment for demonstrations execute the following command from the go-test-bench directory.

```bash
./demo.sh reset
```

This will re-build your test bench with the most recent copies of the go-agent, contrast-service and test bench.

## Troubleshooting

### I updated my config but nothing changed

If you've updated your config locally you will need to reset your environment using the following command from your go-test-bench directory.

```bash
./demo.sh reset
```
