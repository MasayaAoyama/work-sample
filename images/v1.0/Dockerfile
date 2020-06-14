FROM golang:1.13

WORKDIR /go/src/app
COPY ./main.go .

RUN go get -d -v github.com/go-sql-driver/mysql
RUN go get -d -v github.com/Pallinder/go-randomdata

RUN go build -o /app ./main.go

CMD ["/app"]
