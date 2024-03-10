FROM golang:latest

WORKDIR /go/src

RUN go mod init TravelGachaGo

RUN go get github.com/gin-gonic/gin