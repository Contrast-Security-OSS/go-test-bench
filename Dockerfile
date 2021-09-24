FROM golang:1.16 AS builder
ARG FRAMEWORK=std

WORKDIR /build
COPY . .

RUN go mod download
RUN go build \
      -ldflags='-w -s -extldflags "-static"' \
      -o go-test-bench \
      ./cmd/${FRAMEWORK}

FROM scratch
WORKDIR /

COPY --from=builder /build/views /views
COPY --from=builder /build/public /public
COPY --from=builder /build/go-test-bench /go-test-bench
# Copy /etc/passwd for easy path-traversal
COPY --from=builder /etc/passwd /etc/passwd

EXPOSE 8080

ENTRYPOINT ["/go-test-bench"]