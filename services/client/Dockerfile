FROM golang:1.13


WORKDIR $GOPATH/src/github.com/fredericpinaud/client

COPY . .

RUN ls -la

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["client"]