FROM golang:latest

WORKDIR /go/src

RUN go mod init TravelGachaGo

RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jinzhu/gorm
