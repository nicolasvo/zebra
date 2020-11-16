FROM golang:alpine AS build
WORKDIR /go/src/zebra
COPY go.mod go.sum main.go Makefile .
RUN apk add make
RUN go get
ENTRYPOINT ["/bin/sh"]
