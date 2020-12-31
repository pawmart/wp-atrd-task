FROM golang:1.15.6-alpine3.12
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o /app/secretsgcd /app/cmd/secretsgcd
CMD ["/app/secretsgcd"]
