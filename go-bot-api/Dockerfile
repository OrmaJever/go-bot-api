FROM golang:1.20

WORKDIR /usr/go-app

ADD . ./

RUN go mod download
RUN go build -o main .

CMD ["/usr/go-app/main"]
