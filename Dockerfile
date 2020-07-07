FROM golang:1.14 AS builder

# Set the working directory and copy the code over
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

# Fetch dependencies.
RUN go get -d -v

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/app .

# Create the atomizer container
FROM alpine:latest

WORKDIR /

# Copy the atomizer agent to thee new scratch container
COPY ./views /views
COPY ./public /public
COPY --from=builder /go/bin/app /app

EXPOSE 8080

# Execute the atomizer agent
ENTRYPOINT ["./app"]