# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.15 as builder

# Create the /app and /config directories.
#RUN mkdir /app && mkdir /config

# Copy the config.json.template file to /config/config.json
COPY ./config/config.json.template /config/config.json

# Set the working directory to /app.
WORKDIR /app

# Copy all the files from the current directory to /app.
COPY . /app

# Build the binary.
# -mod=readonly ensures that the go.mod and go.sum files are not updated.
RUN go build -mod=readonly -o gramarr ./cmd/gramarr

# Set the default command to run the gramarr binary.
CMD ["/app/gramarr", "-configDir=/config"]
