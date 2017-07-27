FROM golang:1.8-alpine

ENV PWD $GOPATH/src/github.com/stone-payments/gcd

RUN mkdir -p $PWD

WORKDIR $PWD

ADD . ./

RUN go install

CMD ["gcd"]