FROM golang:1.15 as builder

# Set the working directory to /app
WORKDIR /app

# Copy all the files from the current directory to /go/src/app
COPY . . 

# Build the binary 
RUN go build -mod=readonly ./cmd/gramarr

# Use the official golang image as the second stage
FROM golang:1.15

# Create the /app and /config directories
RUN mkdir /app && mkdir /config

# Copy the config.json.template file to /config/config.json
COPY ./config/config.json.template /config/config.json

# Copy the gramarr binary from the builder stage
COPY --from=builder /app/gramarr /app/gramarr

# Set the default command to run the gramarr binary
CMD ["/app/gramarr", "-configDir=/config"]
