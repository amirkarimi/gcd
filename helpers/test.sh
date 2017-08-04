export GCD_TEST_AMOUNT_CONTAINERS=10
for image in "${images[@]}"
do
    echo "Pull "$image "-" $(docker pull $image)
done

for container in `eval echo {1..$GCD_TEST_AMOUNT_CONTAINERS}`
do
    echo "Up c"$container "-" $(docker run -d --name "c"$container alpine sh -c "while true; do sleep 10000; done")
done
go test ./worker -run TestListContainers -cover -v

docker rm -f $(docker ps -aq)
