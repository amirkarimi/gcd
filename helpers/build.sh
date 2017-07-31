if which docker >/dev/null ; then
    mkdir -p ./tmp/build
    docker run --rm \
        --user $(id -u):$(id -g) \
        -e GOPATH=$GOPATH \
        -e USERNAME=$USERNAME \
        -e CGO_ENABLED=0 \
        -v $(pwd):$(pwd) \
        -w $(pwd) \
    golang:1.8-onbuild go build
    mv ./gcd ./tmp/build/
    docker build -t gcd .
   rm -rf ./tmp
    exit 0;
fi

echo "Docker not installed"
