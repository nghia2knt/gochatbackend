FROM golang:1.19-alpine

WORKDIR /go/src/app

COPY . .

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]