FROM golang:1.13

WORKDIR $GOPATH/src/github.com/fredericpinaud/client

COPY client.go .

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["client"]