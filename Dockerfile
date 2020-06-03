FROM golang:1.14 AS build
WORKDIR /go/src
COPY . .

ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -o excommerce .

FROM scratch AS runtime
COPY --from=build /go/src/excommerce ./
EXPOSE 8080/tcp
ENTRYPOINT ["./excommerce"]
