APP_NAME=spiffe-demo-app
DOCKER_IMAGE=spiffe-demo-app
VERSION=latest


.PHONY: all build image clean run

all: build image

build:
	go build -o ./bin/$(APP_NAME)

image:
	# env variables and .ko.yaml are ignored for a some reason
	# so we set them up here manually
	ko build --platform linux/arm64 --preserve-import-paths --local .

run:
	KO_DATA_PATH=./kodata SPIFFE_ENDPOINT_SOCKET=unix:///tmp/spirl/spiffe.sock go run main.go

clean:
	rm -rf ./bin
