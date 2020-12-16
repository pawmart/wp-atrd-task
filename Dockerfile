FROM golang:1 AS build
WORKDIR /src/
COPY ./ /src/
RUN CGO_ENABLED=0 go build -o /secret-server github.com/systemz/wp-atrd-task/cmd/server

FROM scratch
COPY --from=build /secret-server /bin/secret-server
COPY api /bin/api
ENTRYPOINT ["/bin/secret-server"]
