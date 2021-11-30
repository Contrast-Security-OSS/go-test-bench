FROM golang:1.17 AS builder
ARG FRAMEWORK=std

WORKDIR /build
COPY . .

RUN go mod download
RUN go build \
      -ldflags='-w -s -extldflags "-static"' \
      -o go-test-bench \
      ./cmd/${FRAMEWORK}

# Create /etc/passwd to help showcase path traversal vulnerability.
RUN echo "root:x:0:0:root:/root:/bin/bash" >> ./fakepasswd && \
      echo "catty:x:42:42:boh:/etc/contrastsecurity:/bin/bash" >> ./fakepasswd

FROM scratch
WORKDIR /

COPY --from=builder /build/views /views
COPY --from=builder /build/public /public
COPY --from=builder /build/go-test-bench /go-test-bench
COPY --from=builder /build/fakepasswd /etc/passwd

EXPOSE 8080

# default listen address of localhost:8080 does not work
ENTRYPOINT ["/go-test-bench", "-addr=:8080"]