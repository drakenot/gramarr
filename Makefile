.PHONY: binary image run-image

all: image

binary:
		CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo
			
image: binary
		docker build -t gramarr:latest .
			
run-image:
		docker run gramarr:latest
