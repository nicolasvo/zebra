FROM golang:alpine AS build
WORKDIR /go/src/zebra
COPY go.mod go.sum main.go Makefile .
RUN apk add make
RUN make clean
RUN go get
RUN make target=arm build

FROM alpine:latest AS run
WORKDIR /app
COPY . .
COPY --from=build /go/src/zebra/zebra .
CMD ["sh", "run.sh"]