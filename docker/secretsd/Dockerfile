FROM golang:1.15.6-alpine3.12
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o /app/secretsd /app/cmd/secretsd
CMD ["/app/secretsd"]

