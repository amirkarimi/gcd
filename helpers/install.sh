if which go >/dev/null ; then
    go install -v
    echo "gcd installed in $GOPATH/bin"
    exit 0;
fi

echo "Go not installed"