# --- Stage 1: Build the binary ---
FROM golang:1.25-alpine AS builder

# Install git if your private dependencies require it
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy dependency manifests first for efficient layer caching
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY server/main.go ./
COPY server/internal ./internal

# Build the statically linked binary optimized for production
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# --- Stage 2: Minimal runtime image ---
FROM alpine:3.20

# ca-certificates so TLS to Supabase/Postgres works
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/main ./main

EXPOSE 8080
ENTRYPOINT ["./main"]
