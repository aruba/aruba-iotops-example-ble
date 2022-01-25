DOCKER_IMG=aruba-iotops-example-ble:$(shell head ../VERSION)

ifeq ($(shell uname),Darwin)
DOCKER_FLAGS +=--platform=linux/amd64
endif

.PHONY: test
test:
	go test -v ./application

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w -s'  -o ./app ./application

.PHONY: docker
docker: build
	docker build ${DOCKER_FLAGS}  -f ./Dockerfile -t $(DOCKER_IMG) .
