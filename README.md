# gcd

> WIP (Work In Progress)

## Install
```bash
> git pull https://github.com/stone-payments/gcd.git
> make build
```

## To use
```bash
> docker images

REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
gdc                 latest              b5e360ce4d15        2 seconds ago       270 MB
docker.io/golang    1.8-alpine          310e63753884        4 weeks ago         257 MB

> docker run --rm -v /var/run/docker.sock:/var/run/docker.sock gcd

INFO Time: 1501254532651998785
INFO Host: unix:///var/run/docker.sock
INFO Containers total: 1
INFO Images total: 2
INFO Action to containers finished
INFO Action to images finished

INFO Time: 1501263323571366035
INFO Host: unix:///var/run/docker.sock
INFO Containers total: 2
INFO Images total: 3
INFO Container 3ccc041793553a28da25c185afb1a93b270d58c83627602a35a44a5efa683b3a removed successful
INFO Action to containers finished
INFO Image sha256:6833171fe0ad8f221a1f9c099ffcc40fbab6dfb1e70589975cac3355cf08c118 removed successful
INFO Action to images finished

INFO Time: 1501254533669690741
INFO Host: unix:///var/run/docker.sock
INFO Containers total: 1
INFO Images total: 2
INFO Action to containers finished
INFO Action to images finished

```

## Roadmap

- [x] Create logic core
- [x] Create Dockerfile
- [ ] Create tests
- [x] Create _How to use_
- [ ] Create benchmark
