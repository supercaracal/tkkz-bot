SHELL    := /bin/bash
APP_NAME := tkkz-bot

all: build test lint

build:
	go build -ldflags="-s -w" -trimpath -tags timetzdata -o ${APP_NAME}

test:
	go test

lint:
	go vet
	golint -set_exit_status

clean:
	@rm -f ${APP_NAME} main

build-image:
	@docker build -t ${APP_NAME} .
	@docker image prune -f

lint-image:
	@docker run --rm -i hadolint/hadolint < Dockerfile

run-container:
	@docker run --env-file=.env --rm ${APP_NAME}

clean-image:
	@docker rmi -f ${APP_NAME}

.PHONY: all build test lint clean build-image lint-image run-container clean-image
