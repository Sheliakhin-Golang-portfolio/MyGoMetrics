# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary with CGO disabled
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -ldflags="-w -s" -o /mygometrics ./cmd

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary from builder stage
COPY --from=builder /mygometrics /mygometrics

# Expose default port
EXPOSE 9000

# Run as non-root user (distroless images use nonroot user by default)
USER nonroot:nonroot

# Entrypoint
ENTRYPOINT ["/mygometrics"]
