FROM alpine:latest
# packages for TLS openssl ca-certificates
# as an alternative use CGO=0 during building
RUN apk --no-cache add musl-dev
ADD secret-server /
RUN chmod +x /secret-server
ENTRYPOINT /secret-server
