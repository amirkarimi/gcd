all:
	make test
	make bench
	make build
	make push

test:
	@sh ./helpers/test.sh

bench:
	@sh ./helpers/bench.sh

cover: 
	@sh ./helpers/cover.sh

build:
	@sh ./helpers/build.sh

push:
	@sh ./helpers/push.sh