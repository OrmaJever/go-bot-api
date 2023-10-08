FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED 0
ENV GOOS linux
RUN go build -o main ./cmd

FROM alpine
WORKDIR /app
RUN apk add --update --no-cache ca-certificates
COPY --from=builder /usr/share/zoneinfo/Europe/Kiev /usr/share/zoneinfo/Europe/Kiev
ENV TZ Europe/Kiev
COPY . .
COPY --from=builder /app/main /app/main
ENTRYPOINT ["./main"]