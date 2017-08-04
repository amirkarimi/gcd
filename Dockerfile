FROM alpine

ENV GCD_DOCKER_HOST "/var/run/docker.sock"
ENV GCD_SWEEP_INTERVAL "1"
ENV GCD_DOCKER_API_VERSION "1.24"

ADD ./tmp/build/* ./usr/bin

CMD ["gcd"]