FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o auth .

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/auth .
EXPOSE 8080
CMD ["./auth"]