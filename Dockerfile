FROM golang:1.10

WORKDIR /go/cmd/main
COPY ./cmd/main .
CMD ["/go/cmd/main/run.sh"]