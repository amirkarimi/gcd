if which docker >/dev/null ; then
    docker build -t gcd .
    exit 0;
fi

echo "Docker not installed"