FROM golang:1.13 AS build

RUN mkdir -p /go/src/github.com/alcmoraes/gramarr

WORKDIR /go/src/github.com/alcmoraes/gramarr

COPY . .

RUN go get

RUN mkdir -p /app

RUN go build -o /app/gramarr

COPY config.json /app/config.json

CMD ["/app/gramarr", "-configDir=/app"]