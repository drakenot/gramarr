FROM golang:1.13 AS build

RUN mkdir -p /go/src/github.com/drakenot/gramarr

WORKDIR /go/src/github.com/drakenot/gramarr

COPY . .

RUN go get

RUN mkdir -p /app

RUN mkdir -p /config

RUN go build -o /app/gramarr

COPY config.json.template /config/config.json

CMD ["/app/gramarr", "-configDir=/config"]
