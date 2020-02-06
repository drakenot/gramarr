FROM golang:onbuild
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get
RUN go build -o gramarr .
CMD ["/app/gramarr"]
