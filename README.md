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

REPOSITORY               TAG                 IMAGE ID            CREATED             SIZE
docker.io/guiferpa/gcd   latest              04ba50851638        16 seconds ago      9.7 MB
gcd                      latest              04ba50851638        16 seconds ago      9.7 MB
guiferpa/gcd             latest              04ba50851638        16 seconds ago      9.7 MB
docker.io/golang         1.8-onbuild         5d82e356477f        2 weeks ago         699 MB
docker.io/alpine         latest              7328f6f8b418        6 weeks ago         3.97 MB

> docker run --name gcd -v /var/run/docker.sock:/var/run/docker.sock gcd

INFO Time: 1501884318075906011
INFO State: running
INFO Docker API Version: v1.24

INFO Time: 1501884319082624665
INFO Host: /var/run/docker.sock
INFO Containers total: 1
INFO Images total: 2
INFO Action to containers finished
OK Image sha256:cb768d4c340233275d8dab6160de29c849b5b5dea63bce09a1c683205a75d855 removed successful
OK Image sha256:5d82e356477f55c857e52b672d244e4d72f065c81cdd5010efdad6eb26ea5088 removed successful
INFO Action to images finished

INFO Time: 1502305011634686589
INFO Host: unix:///var/run/docker.sock
INFO Containers total: 1
INFO Images total: 1
SKIP Container 885c3d698d3edf5c2dc7fb21e94f223bd3c28166b2a39f19867e2a3920675246 skipped, Status: Exited (0) 15 minutes ago
INFO Action to containers finished
SKIP Image sha256:b39de1d79a536613ed278b5a4d500bcdefa74ecdcecc3a8d655271c3d5bf813f
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
- [x] Create tests (_work in progress_)
- [x] Create _How to use_
- [x] Create benchmark (_work in progress_)
