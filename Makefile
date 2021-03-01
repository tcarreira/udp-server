DOCKER_IMAGE_BASE=tcarreira/udp-server
DOCKER_TAG=$(shell date +%Y%m%d_%H%M)

build:
	go build -i -v -o udp-server

docker-build:
	docker build -t ${DOCKER_IMAGE_BASE}:latest .
	chmod +x udp-server
