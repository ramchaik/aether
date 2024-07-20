FROM golang:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

FROM alpine:latest  
ENV PATH="/bin:${PATH}"
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api .

EXPOSE 8080
CMD ["./api"]
