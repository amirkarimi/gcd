FROM alpine

ADD ./tmp/build/* ./usr/bin

CMD ["gcd"]