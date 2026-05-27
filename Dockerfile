# Stage 1: Build the application
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o red-engine .

# Stage 2: Construct the bare execution container
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -S redgroup && adduser -S reduser -G redgroup

WORKDIR /app
COPY --from=builder /app/red-engine .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Ensure the user owns the directory
RUN chown -R reduser:redgroup /app

# Switch to the restricted user
USER reduser

EXPOSE 8080
VOLUME ["/app/data"]
CMD ["./red-engine"]
