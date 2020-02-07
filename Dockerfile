FROM golang:1.13 as builder

WORKDIR /go/src/github.com/gramarr
COPY "${PWD}" /go/src/github.com/gramarr
RUN go get && go build -o ./build/gramarr .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/gramarr/build/gramarr /usr/bin/gramarr
