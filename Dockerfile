# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.15 as builder

# Copy local code to the container image.
WORKDIR /go/src/app
COPY . .

# Build the binary.
# -mod=readonly ensures that the go.mod and go.sum files are not updated.
RUN go build -mod=readonly -o gramarr

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/app/gramarr /app/gramarr

# Copy the config.json.template file to /config/config.json
COPY config.json.template /config/config.json

CMD ["/app/grammar", "-configDir=/config"]
