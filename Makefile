all:
	make test
	make bench
	make build

test:
	@sh ./helpers/test.sh

bench:
	@sh ./helpers/bench.sh

cover: 
	@sh ./helpers/cover.sh

build:
	@sh ./helpers/build.sh
