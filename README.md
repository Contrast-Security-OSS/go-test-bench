# Go Test Bench

[![CI](https://github.com/Contrast-Security-OSS/go-test-bench/workflows/CI/badge.svg)](https://github.com/Contrast-Security-OSS/go-test-bench/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Contrast-Security-OSS/go-test-bench)](https://goreportcard.com/report/github.com/Contrast-Security-OSS/go-test-bench)
[![GoDoc](https://godoc.org/github.com/Contrast-Security-OSS/go-test-bench?status.svg)](https://pkg.go.dev/github.com/Contrast-Security-OSS/go-test-bench)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Intentionally vulnerable go app. Used Go's standard library, `net/http`,
for client/server implementations. For more info on this framework, visit
[net/http](https://golang.org/pkg/net/http/).

The go-test-bench application includes vulnerabilities from the OWASP Top
10 and is intended to be used as an educational tool for developers and
security professionals. Any maintainers are welcome to make pull requests.

> For customer demonstrations [click here to follow the Demo.md readme](./Demo.md).

## Installation Requirements

- [Go 1.16 or higher](https://golang.org/dl/)

- *Optional* [Docker for Mac](https://www.docker.com/docker-mac)

## How to Run Locally

* standard library
```bash
    go build -o app ./cmd/std
    ./app
```

To run with gin instead, substitute `gin` for `std` in the build command; likewise for `go-swagger` and `julienschmidt`.

View app at [http://localhost:8080](http://localhost:8080)

## How to Run Using Docker

```bash
    # To stand up application with standard library
    docker-compose up std -d
    # Alternative: To stand up application with Gin framework
    # docker-compose up gin -d
    # To stand up application with Julienschmidt/httprouter:
    # docker-compose up julienschmidt -d

    # To stop app container and related service containers
    docker-compose down
```

View app at [http://0.0.0.0:8080](http://0.0.0.0:8080)

## Acknowledgements

The development [team](docs/acknowledgements.md).


## Experimenting with the code

### organization

* code for vulnerable functions is located in `internal/`
  * exception: vulnerable functions from a particular framework (see below)
* framework-specific code is located under `cmd/` and `pkg/`
* html templates and css are under `views/`
* vulnerability and route data is in go structs, 
  located in the relevant package under `internal/`

### quirks

Each framework is different. We've tried to separate framework logic from vulnerability logic so that adding a framework necessitates a minimum of changes to vulnerability logic, and vice versa.

#### swagger
Swagger is a bit unique, in that it has a lot of generated code and requires a swagger spec. To maintain a single source of truth, we generate the swagger spec from our route data. We also generate boilerplate tying a route handler to each swagger endpoint.

For details, see [cmd/go-swagger/README.md](cmd/go-swagger/README.md)
