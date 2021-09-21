FROM golang:1.16 AS builder
ARG FRAMEWORK=std

COPY . /build
WORKDIR /build/cmd/${FRAMEWORK}

RUN go mod download
RUN go build \
      -ldflags='-w -s -extldflags "-static"' \
      -o /build/go-test-bench

FROM scratch
WORKDIR /

COPY --from=builder /build/views /views
COPY --from=builder /build/public /public
COPY --from=builder /build/go-test-bench /go-test-bench

EXPOSE 8080

ENTRYPOINT ["/go-test-bench"]