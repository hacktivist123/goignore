# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /go/src/goignore

# Install git
RUN apk add --no-cache git

# Clone the repository and build the goignore executable
RUN git clone https://github.com/hacktivist123/goignore . && \
  go build -o /go/bin/goignore ./cmd/goignore

# goignore Container Image
FROM alpine:latest

# Maintainer info
LABEL org.opencontainers.image.authors="Shedrack Akintayo" \
  org.opencontainers.image.description="Container image for https://github.com/hacktivist123/goignore"

# Set the working directory
WORKDIR /goignore

# Copy the goignore executable from the builder stage
COPY --from=builder /go/bin/goignore /usr/local/bin/goignore

# Set the entrypoint
ENTRYPOINT [ "goignore" ]
