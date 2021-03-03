FROM golang:1.16-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o /secret-server ./cmd/secret-server

CMD ["/secret-server"]