# gcd

[![Build Status](https://travis-ci.org/stone-payments/gcd.svg?branch=master)](https://travis-ci.org/stone-payments/gcd)

## Install

```bash
> git pull https://github.com/stone-payments/gcd.git
> make build
```

## To use

### Short command
```bash
> docker run --name gcd -v /var/run/docker.sock:/var/run/docker.sock guiferpa/gcd
```
### Running builded image
```bash
> docker images

REPOSITORY               TAG                 IMAGE ID            CREATED             SIZE
guiferpa/gcd             latest              04ba50851638        16 seconds ago      9.7 MB
docker.io/golang         1.8-onbuild         5d82e356477f        2 weeks ago         699 MB
docker.io/alpine         latest              7328f6f8b418        6 weeks ago         3.97 MB

> docker run --name gcd -v /var/run/docker.sock:/var/run/docker.sock gcd
```

## Flags

- __-target__ Set docker host target
- __-sweep-interval__ Set interval between sweep
- __-remove-images__ Set enable to remove any image
- __-remove-healthy-container__ Set enable to remove healthy container

## Docker environment configuration

- __GCD_DOCKER_HOST:__ Set your path for docker.sock, by default use `/var/run/docker.sock:/var/run/docker.sock`
- __GCD_SWEEP_INTERVAL:__ Set your interval to sweep, by default use 1 second
- __GCD_REMOVE_IMAGES__: Set `true` or `false` to remove images, by default use `true`
- __GCD_REMOVE_HEALTHY_CONTAINERS_EXITED__: Remove containers with exited code equal 0, by default use `false`
