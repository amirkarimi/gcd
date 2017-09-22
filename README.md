# gcd

[![Build Status](https://travis-ci.org/stone-payments/gcd.svg?branch=master)](https://travis-ci.org/stone-payments/gcd)

## Description
This project is a garbage collector for docker images and container, it exists because today our machines contains very garbage from Docker components.

## Install

```bash
> git clone https://github.com/stone-payments/gcd.git
> make build
```

## To use

### Running builded binary
```bash
> ./bin/gcd
```

### Use docker image with short command
```bash
> docker run --name gcd -v /var/run/docker.sock:/var/run/docker.sock guiferpa/gcd
```

### Building docker image
```bash
> make build-image
```
> :warning: This project use multi-stage build to build docker image, required use >=17.05 docker version

### Running builded image
```bash
> docker images

REPOSITORY               TAG                 IMAGE ID            CREATED             SIZE
guiferpa/gcd             latest              04ba50851638        16 seconds ago      9.7 MB
docker.io/golang         1.8-onbuild         5d82e356477f        2 weeks ago         699 MB
docker.io/alpine         latest              7328f6f8b418        6 weeks ago         3.97 MB

> docker run --name gcd -v /var/run/docker.sock:/var/run/docker.sock guiferpa/gcd
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
