FROM golang:alpine
WORKDIR /go/src/zebra
COPY . .
RUN apk add make
RUN go get
ENTRYPOINT ["/bin/sh"]
