FROM golang:1.16-alpine

COPY . /go/src/app
WORKDIR /go/src/app

RUN go mod download

RUN go build -o api.go

EXPOSE 8010

CMD ["./api.go"]
