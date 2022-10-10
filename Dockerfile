FROM golang:1.19

RUN go install github.com/registrobr/rdap-client@latest

ENTRYPOINT ["/go/bin/rdap-client"]
