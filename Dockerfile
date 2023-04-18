# Using a multi-stage build to avoid compilation in the Docker image
FROM golang:1.18 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o port-api cmd/port-api/main.go

# Use a minimal image for running the service
FROM alpine:3.14

# Address security concerns by creating a non-root user to run the service
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set the working directory for the application
WORKDIR /app

COPY --from=builder /app/port-api /app/port-api
COPY --from=builder /app/ports.json /app/ports.json

CMD ["/app/port-api"]

EXPOSE 8080
