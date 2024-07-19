FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

FROM docker:dind

WORKDIR /app

COPY --from=builder /app/worker .
COPY --from=builder /app/secure-build.dockerfile .
COPY .env .

ENV DOCKER_HOST=unix:///var/run/docker.sock

EXPOSE 8080

COPY <<EOF /start.sh
#!/bin/sh

# Start dockerd with cgroup v1
dockerd --exec-opt native.cgroupdriver=cgroupfs &

# Wait for Docker to become available
timeout=60
while ! docker info >/dev/null 2>&1; do
    timeout=$((timeout - 1))
    if [ $timeout -le 0 ]; then
        echo "Timed out waiting for Docker to start"
        exit 1
    fi
    sleep 1
done

echo "Docker is available, starting worker..."
exec ./worker
EOF

RUN chmod +x /start.sh

ENTRYPOINT ["/start.sh"]