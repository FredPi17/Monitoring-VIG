FROM golang:1.13

WORKDIR $GOPATH/src/github.com/fredericpinaud/server

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["server"]