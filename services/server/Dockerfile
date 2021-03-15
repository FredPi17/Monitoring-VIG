FROM golang:1.13

WORKDIR $GOPATH/src/github.com/fredericpinaud/server

COPY server.go .

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["server"]