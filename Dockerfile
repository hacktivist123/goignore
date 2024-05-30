# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /go/src/goignore

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code into the container
COPY . .

# Build the goignore executable
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/goignore ./cmd/goignore

# Final stage
FROM scratch

# Maintainer info
LABEL org.opencontainers.image.authors="Shedrack Akintayo" \
  org.opencontainers.image.description="Container image for https://github.com/hacktivist123/goignore"

# Copy the goignore executable from the builder stage
COPY --from=builder /go/bin/goignore /usr/local/bin/goignore

# Set the entrypoint
ENTRYPOINT [ "/usr/local/bin/goignore" ]
