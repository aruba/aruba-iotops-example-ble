#Configure your docker image name with registry url
DOCKER_IMG=aruba-iotops-example-ble:$(shell head ../VERSION)
imagename?=images

.PHONY: build
build:
	go build -o app ./application

.PHONY: test
test:
	go test -v ./application

.PHONY: docker-local
docker-local:
	docker buildx build --load -t $(DOCKER_IMG) .

#Support both data collector and AP as connector
IMAGE_PLATFORM +=--platform=linux/amd64,linux/arm64,linux/arm/v7
.PHONY: publish-docker
publish-docker:
	docker buildx build ${IMAGE_PLATFORM} --push -t $(DOCKER_IMG) .

#Build a tar to upload in Aruba Developer Portal
.PHONY: image-tar
image-tar:
	docker buildx build ${IMAGE_PLATFORM} -o type=oci,dest=$(imagename).tar .
