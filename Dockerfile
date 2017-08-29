FROM golang:1.8 as builder

ENV WD /go/src/github.com/stone-payments/gcd

COPY . $WD

WORKDIR $WD

RUN CGO_ENABLED=0 go build

FROM alpine

ENV WD /go/src/github.com/stone-payments/gcd

COPY --from=builder $WD/gcd .

ENV GCD_DOCKER_HOST "/var/run/docker.sock"
ENV GCD_SWEEP_INTERVAL "1"
ENV GCD_DOCKER_API_VERSION "1.24"
ENV GCD_REMOVE_IMAGES "true"
ENV GCD_REMOVE_CONTAINERS_EXITED "false"

CMD ["./gcd"]
