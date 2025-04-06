# Use a specific version of Golang for better reproducibility
FROM golang:1.24.2-alpine3.19 as builder


# Set the working directory inside the container
WORKDIR /app

# Copy only the dependency files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies with specific version
RUN go mod download -x

# Copy the source code into the container
COPY . .

# Build the Go application with optimizations and security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE} -X main.Commit=${COMMIT_SHA}" \
    -o main ./cmd/main.go

# Use a minimal base image for the final container
FROM public.ecr.aws/lambda/provided:al2023

# Add non-root user for security
RUN adduser -D -g '' appuser

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder --chown=appuser:appuser /app/main .

# Switch to non-root user
USER appuser

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Add metadata
LABEL org.opencontainers.image.version=${VERSION} \
      org.opencontainers.image.created=${BUILD_DATE} \
      org.opencontainers.image.revision=${COMMIT_SHA} \
      org.opencontainers.image.title="Iris Application" \
      org.opencontainers.image.description="Iris Application Container"

# Command to run the application
ENTRYPOINT ["./main"]