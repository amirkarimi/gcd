# gcd

[![Build Status](https://travis-ci.org/stone-payments/gcd.svg?branch=master)](https://travis-ci.org/stone-payments/gcd)
[![Coverage Status](https://coveralls.io/repos/github/stone-payments/gcd/badge.svg?branch=master)](https://coveralls.io/github/stone-payments/gcd?branch=master)

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

REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
gdc                 latest              b5e360ce4d15        2 seconds ago       270 MB
docker.io/golang    1.8-alpine          310e63753884        4 weeks ago         257 MB

> docker run --name gcd -v /var/run/docker.sock:/var/run/docker.sock gcd

INFO Time: 1501884318075906011
INFO State: running
INFO Docker API Version: v1.24

INFO Time: 1501884319082624665
INFO Host: /var/run/docker.sock
INFO Containers total: 1
INFO Images total: 8
INFO Action to containers finished
OK Image sha256:cb768d4c340233275d8dab6160de29c849b5b5dea63bce09a1c683205a75d855 removed successful
OK Image sha256:5d82e356477f55c857e52b672d244e4d72f065c81cdd5010efdad6eb26ea5088 removed successful
INFO Action to images finished

INFO Time: 1501884320449422524
INFO Host: /var/run/docker.sock
INFO Containers total: 1
INFO Images total: 6
INFO Action to containers finished
INFO Action to images finished

```

## Configuration

- __GCD_DOCKER_HOST:__ Set your path for docker.sock, by default use `/var/run/docker.sock:/var/run/docker.sock`
- __GCD_SWEEP_INTERVAL:__ Set your interval to sweep, by default use 1 second
- __GCD_REMOVE_IMAGES__: Set `true` or `false` to remove images, by default use `true`
- __GCD_REMOVE_CONTAINERS_EXITED__: Remove containers with exited code equal 0, by default use `false`

## Roadmap

- [x] Create logic core
- [x] Create Dockerfile
- [ ] Create tests
- [x] Create _How to use_
- [ ] Create benchmark
