# Stage 1: Build the optimized static binary
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o red-engine .

# Stage 2: Construct the bare execution container
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/red-engine .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expose server port and bind state volumes
EXPOSE 8080
VOLUME ["/root/data"]
CMD ["./red-engine"]
