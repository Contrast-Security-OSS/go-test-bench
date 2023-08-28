ARG GO_VERSION=latest
FROM golang:${GO_VERSION} AS builder

WORKDIR /build

COPY . .

ARG FRAMEWORK=std
RUN CGO_ENABLED=0 go build -o go-test-bench ./cmd/${FRAMEWORK}

# Move the finished build to a more minimal container
FROM scratch
COPY --from=builder /build/views ./views
COPY --from=builder /build/public ./public
COPY --from=builder /build/go-test-bench ./go-test-bench

ENTRYPOINT ["./go-test-bench", "-addr=:8080"]
