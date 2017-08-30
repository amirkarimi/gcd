all:
	@make test
	@make bench
	@make cover
	@make build

docker:
	@make docker-build
	@make docker-up

test:
	@sh ./helpers/test.sh

bench:
	@sh ./helpers/bench.sh

cover: 
	@sh ./helpers/cover.sh

build:
	@sh ./helpers/build.sh

docker-build:
	docker build -t gcd .

docker-up:
	docker run -v $$DOCKER_SOCK:/var/run/docker.sock gcd
