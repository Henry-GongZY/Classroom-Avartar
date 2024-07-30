FROM golang:latest
MAINTAINER Henry-GongZY "1811358@mail.nankai.edu.cn"

WORKDIR $GOPATH/src/
ADD . $GOPATH/src/

RUN go build -o Server ./Server.go

EXPOSE 5665

ENTRYPOINT ["./Server"]