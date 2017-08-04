all:
	make test
	make build
	make push

test:
	@sh ./helpers/test.sh

cover: 
	@sh ./helpers/cover.sh

build:
	@sh ./helpers/build.sh

push:
	@sh ./helpers/push.sh