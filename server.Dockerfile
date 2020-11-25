FROM golang:1.14.6-alpine3.12 as builder

COPY go.mod go.sum /go/src/github.com/maciejem/secret/
WORKDIR /go/src/github.com/maciejem/secret
RUN go mod download
COPY . /go/src/github.com/maciejem/secret
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/secret github.com/maciejem/secret


FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
RUN apk update && apk add bash

COPY --from=builder /go/src/github.com/maciejem/secret/build/secret /usr/bin/secret

COPY --from=builder /go/src/github.com/maciejem/secret/scripts scripts 
RUN chmod +x scripts/wait-for-it.sh

EXPOSE 8080 8080

ENTRYPOINT [ "/bin/sh", "-c" ]
