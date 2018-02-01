# Build binary
FROM golang:1.8 as builder
WORKDIR /go/src/github.com/drakenot/gramarr/
COPY . . 
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Create minimal image with binary 
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/drakenot/gramarr/app .
ENTRYPOINT ["/app", "--config=/config"]
VOLUME ["/config"]
