FROM golang:1.20

WORKDIR /usr/go-app

ADD models models
ADD packages packages
ADD services services
ADD telegram telegram
ADD .env .env
ADD go.mod go.mod
ADD go.sum go.sum
ADD handler.go handler.go
ADD main.go main.go


RUN go mod download
RUN go build -o main .

CMD ["/usr/go-app/main"]
