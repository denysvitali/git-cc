# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for testing
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE} -X main.builtBy=docker" \
    -a -installsuffix cgo -o git-cc .

# Test the build
RUN ./git-cc --version

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates git

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/git-cc .

# Add entrypoint
ENTRYPOINT ["./git-cc"]
CMD ["--version"]