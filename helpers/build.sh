if which docker >/dev/null ; then
    mkdir -p ./tmp/build
    docker run --rm \
        -e GOPATH=$GOPATH \
        -e CGO_ENABLED=0 \
        -v $(pwd):$(pwd):Z \
        -w $(pwd) \
        -it \
        golang:1.8-onbuild go build
    mv ./gcd ./tmp/build/
    docker build -t gcd .
    docker tag gcd guiferpa/gcd
    docker tag gcd docker.io/guiferpa/gcd
    rm -rf ./tmp
    exit 0;
fi

echo "Docker not installed"
