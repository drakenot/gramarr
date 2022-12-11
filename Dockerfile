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

# Copy the binary to the /app directory.
RUN mkdir /app && cp gramarr /app/gramarr

# Copy the config.json.template file to /config/config.json
COPY config.json.template /config/config.json

# Set the default command to run the gramarr binary.
CMD ["/app/gramarr", "-configDir=/config"]
