.PHONY: build install

build:
	docker build -t Docker.build

install:
	docker run --rm --entrypoint cat mtp-util /usr/local/bin/mtp > mtp
	chmod +x mtp
