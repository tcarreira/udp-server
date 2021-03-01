DOCKER_IMAGE_BASE=tcarreira/udp-server
DOCKER_TAG=$(shell date +%Y%m%d_%H%M)

build:
	go build -i -v -o udp-server

docker-build:
	docker build -t ${DOCKER_IMAGE_BASE}:latest .
	chmod +x udp-server

docker-publish:
	docker tag ${DOCKER_IMAGE_BASE}:latest ${DOCKER_IMAGE_BASE}:${DOCKER_TAG}
	docker push ${DOCKER_IMAGE_BASE}:latest
	docker push ${DOCKER_IMAGE_BASE}:${DOCKER_TAG}

fmt:
	go fmt

