# Go Test Bench

[![CI](https://github.com/Contrast-Security-OSS/go-test-bench/workflows/CI/badge.svg)](https://github.com/Contrast-Security-OSS/go-test-bench/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Contrast-Security-OSS/go-test-bench)](https://goreportcard.com/report/github.com/Contrast-Security-OSS/go-test-bench)
[![GoDoc](https://godoc.org/github.com/Contrast-Security-OSS/go-test-bench?status.svg)](https://pkg.go.dev/github.com/Contrast-Security-OSS/go-test-bench)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Intentionally vulnerable go app. Used Go's standard library, `net/http`,
for client/server implentations. For more info on this framework, visit
[net/http](https://golang.org/pkg/net/http/).

The go-test-bench application includes vulnerabilities from the OWASP Top
10 and is intended to be used as an educational tool for developers and
security professionals. Any maintainers are welcome to make pull requests.

## Installation/Requirements

- On a mac, first install go (>=1.13).

- Install [Docker for Mac](https://www.docker.com/docker-mac)

## How to Run Locally

```bash
    go build app.go
    ./app
```

View app at [http://localhost:8080](http://localhost:8080)

## How to Run Using Docker

```bash
    # To stand up application with mongo
    docker-compose up -d

    # To stop app container and related service containers (ie mongo)
    docker-compose down
```

View app at [http://0.0.0.0:8080](http://0.0.0.0:8080)

## Acknowledgements

The development [team](docs/acknowledgements.md).

