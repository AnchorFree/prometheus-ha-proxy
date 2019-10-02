## binarybuild
##
FROM golang:1.13.1-alpine3.10 as binarybuilder

# Enable support of go modules by default
ENV GO111MODULE on
ENV PROJECT_NAME prom-ha-proxy

# Warming modules cache with project dependencies
WORKDIR /go/src/${PROJECT_NAME}
COPY go.mod go.sum ./
RUN go mod download

# Copy project source code to WORKDIR
COPY . .

# Run tests and build on success
ENV CGO_ENABLED 0
RUN go test -v ./... && GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o /build/${PROJECT_NAME}

## Final container stage
##
FROM alpine
ENV BINARY prom-ha-proxy
WORKDIR /app
COPY --from=binarybuilder /build/${BINARY} bin/${BINARY}

EXPOSE 9090
ENTRYPOINT ["/app/bin/prom-ha-proxy"]
