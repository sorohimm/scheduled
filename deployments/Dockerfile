FROM golang:1.17.1-alpine3.14 as build-env

WORKDIR /schd/cmd
COPY . /schd

RUN go build -o uran .

WORKDIR /
FROM alpine:latest
COPY --from=build-env schd/cmd/uran /schd/cmd/uran

ENTRYPOINT ["schd/cmd/uran"]

