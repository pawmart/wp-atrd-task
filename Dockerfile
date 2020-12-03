FROM golang:1.15 as builder

WORKDIR /goapp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o secretapp cmd/main.go

FROM gcr.io/distroless/base-debian10

COPY --from=builder /goapp/ /

EXPOSE 3000

CMD ["/secretapp"]