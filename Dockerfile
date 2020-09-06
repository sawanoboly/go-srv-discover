FROM golang:alpine

RUN apk add git --no-cache

COPY . /go/src/github.com/sawanoboly/go-srv-discover

WORKDIR /go/src/github.com/sawanoboly/go-srv-discover
RUN go get -v -t -d ./...
CMD ["go", "build", "-o", "build/gsr.alpine", "-v", "."]
