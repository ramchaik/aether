FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# Start a new stage from scratch
FROM alpine:latest  

ENV PATH="/bin:${PATH}"

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/worker .

# Expose port 8080 to the outside world
EXPOSE 8080

CMD ["./worker"]
