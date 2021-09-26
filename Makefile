SHELL      := /bin/bash -euo pipefail
APP_NAME   := tkkz-bot
BRAIN_PORT ?= 3000
GOBIN      ?= $(shell go env GOPATH)/bin

all: build test lint

build: GOOS        ?= $(shell go env GOOS)
build: GOARCH      ?= $(shell go env GOARCH)
build: CGO_ENABLED ?= $(shell go env CGO_ENABLED)
build: FLAGS       += -ldflags="-s -w"
build: FLAGS       += -trimpath
build: FLAGS       += -tags timetzdata
build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build ${FLAGS} -o ${APP_NAME}

test:
	@go clean -testcache
	@go test -race ./...

${GOBIN}/golint:
	go install golang.org/x/lint/golint@latest

lint: ${GOBIN}/golint
	@go vet ./...
	@golint -set_exit_status ./...

clean:
	@rm -f ${APP_NAME} main *.test *.out

run:
	@BRAIN_URL=http://127.0.0.1:${BRAIN_PORT} ./${APP_NAME}

run-debug:
	@BRAIN_URL=http://127.0.0.1:${BRAIN_PORT} ./${APP_NAME} -debug

build-image:
	@docker build -t ${APP_NAME} .

lint-image:
	@docker run --rm -i hadolint/hadolint < Dockerfile

run-container:
	@docker run --env-file=.env --rm ${APP_NAME}

clean-image:
	@docker rmi -f ${APP_NAME}
	@docker image prune -f

run-brain:
	@docker-compose -f docker-compose.development.yml up

stop-brain:
	@docker-compose -f docker-compose.development.yml down

redis-aof:
	@docker-compose exec redis redis-cli bgrewriteaof
