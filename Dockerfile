# Stage 1 - Build binary
FROM golang:1.8 as builder
WORKDIR /go/src/github.com/drakenot/gramarr/
COPY . . 
RUN go get -d -v ./... && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gramarr .

# Stage 2 - Create minimal image with binary 
FROM alpine:latest  
ADD https://github.com/just-containers/s6-overlay/releases/download/v1.21.2.2/s6-overlay-amd64.tar.gz /tmp/
COPY --from=builder /go/src/github.com/drakenot/gramarr/gramarr /app/
COPY config.json.template /app/
COPY docker/root/ /

RUN \
    # install packages
    apk add --update \
	    ca-certificates \
	    shadow && \

    rm -rf /var/cache/apk/* && \

    # install s6-overlay
    tar xzf /tmp/s6-overlay-amd64.tar.gz -C / && \

    # make folders
    mkdir -p /config && \

    # create user
    useradd -u 1000 -U -d /config -s /bin/false gram && \
    usermod -G users gram

VOLUME ["/config"]

ENTRYPOINT ["/init"]
