# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (required for some Go dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
# If go.sum doesn't exist, this will create it
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Generate Swagger docs (required for build since main.go imports docs package)
# Install swag and add to PATH
ENV PATH=$PATH:/go/bin
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go --output docs --parseDependency --parseInternal
# Update go.sum after swagger generation (swag may introduce new dependencies)
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o workmgmt-api \
    main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests, timezone data, and wget for health checks
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/workmgmt-api .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8089

# Set environment variables
ENV PORT=8089
ENV ENVIRONMENT=production
ENV GIN_MODE=release

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8089/api/v1/health || exit 1

# Run the application
CMD ["./workmgmt-api"]
