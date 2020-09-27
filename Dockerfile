# Build
FROM golang:1.14-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .


# Run
FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /root

COPY --from=builder /app/main /root
COPY config/conf.yaml /root/config/conf.yaml

EXPOSE 8080

CMD ["./main"]
