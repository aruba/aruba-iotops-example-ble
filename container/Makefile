#Configure your docker image name with registry url
DOCKER_IMG=aruba-iotops-example-ble:$(shell head ../VERSION)

#Support both data collector and AP as connector
DOCKER_FLAGS +=--platform=linux/amd64,linux/arm64 --push

.PHONY: build
build:
	go build -o app ./application

.PHONY: test
test:
	go test -v ./application

.PHONY: docker
docker:
	docker buildx build ${DOCKER_FLAGS} -t $(DOCKER_IMG) .
