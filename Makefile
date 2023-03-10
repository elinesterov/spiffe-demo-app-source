APP_NAME=spiffe-demo-app
DOCKER_IMAGE=spiffe-demo-app
VERSION=latest


.PHONY: all build image clean

all: build image

build:
	go build -o ./bin/$(APP_NAME)

image:
	# env variables and .ko.yaml are ignored for a some reason
	# so we set them up here manually
	ko build --platform linux/arm64 --preserve-import-paths --local .

clean:
	rm -rf ./bin