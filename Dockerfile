FROM golang:1.14 AS build
WORKDIR /go/src
COPY go ./go
COPY vendor ./vendor
COPY main.go go.mod go.sum ./

ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -o openapi .

FROM scratch AS runtime
COPY --from=build /go/src/openapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]
