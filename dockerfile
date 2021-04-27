
############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /go/src/cjlapao/http-loadtester-go

COPY . .

# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/servicebuscli

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy our static executable.
COPY --from=builder /go/bin/servicebuscli /go/bin/servicebuscli
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV SERVICEBUS_CONNECTION_STRING "Endpoint=sb://sb-defaul-0x0-4172098119.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=8e6d83fddP5ArBRuJJR9Cd4rdhy8bNR77KarlHpUmqs="
ENV SERVICEBUS_CLI_HTTP_PORT 80
EXPOSE 80
ENTRYPOINT ["/go/bin/servicebuscli", "api"]