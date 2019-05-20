FROM golang:1.9


COPY . /go/src/github.com/mishuk-sk/Go-Bank-Transactions
WORKDIR /go/src/github.com/mishuk-sk/Go-Bank-Transactions

RUN go get -v .