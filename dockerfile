############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

ARG TAG_VERSION

# Enable this if you want to have private repositories to be collected
# ARG GITHUB_USER
# ARG GITHUB_TOKEN

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
# Adding the CA certificates
RUN apk --no-cache add ca-certificates
# Adding the sed tool for versioning
RUN apk add sed

# Enable this if you want to have private repositories to be collected
# RUN git config --global url."https://${GITHUB_USER}:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"


WORKDIR /go/src/cjlapao/servicebuscli

COPY . .

WORKDIR /go/src/cjlapao/servicebuscli/src

# Updating the main variable.
RUN sed -i "s/^var ver = \"[[:digit:]]\+\.[[:digit:]]\+\.[[:digit:]]\+\"/var ver = \"${TAG_VERSION}\"/g" main.go

# Build the binary.
RUN GIT_TERMINAL_PROMPT=1 CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/servicebuscli


############################
# STEP 2 build a small image
############################
FROM scratch

# Copy SSL Certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy our static executable.
COPY --from=builder /go/bin/servicebuscli /go/bin/servicebuscli

# Run the project binary.
EXPOSE 5000

ENTRYPOINT ["/go/bin/servicebuscli", "--api"]