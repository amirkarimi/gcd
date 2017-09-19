FROM golang:1.8-onbuild as builder

ENV WD /go/src/github.com/stone-payments/gcd

COPY . $WD

WORKDIR $WD

RUN GOOS=linux CGO_ENABLED=0 go build

FROM alpine

ENV WD /go/src/github.com/stone-payments/gcd

COPY --from=builder $WD/gcd .

ENV GCD_DOCKER_HOST "tcp:///var/run/docker.sock"
ENV GCD_SWEEP_INTERVAL "1"
ENV GCD_REMOVE_IMAGES "false"
ENV GCD_REMOVE_HEALTHY_CONTAINERS_EXITED "true"

CMD ./gcd -target=$GCD_DOCKER_HOST \ 
          -sweep-interval=$GCD_SWEEP_INTERVAL \
          -remove-images=$GCD_REMOVE_IMAGES \
          -remove-healthy-container=$GCD_REMOVE_HEALTHY_CONTAINERS_EXITED
