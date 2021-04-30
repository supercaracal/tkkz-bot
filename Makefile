SHELL      := /bin/bash -e -u -o pipefail
APP_NAME   := tkkz-bot
BRAIN_PORT ?= 3000

all: build test lint

build:
	@go build -ldflags="-s -w" -trimpath -tags timetzdata -o ${APP_NAME}

test:
	@go clean -testcache
	@go test -race ./...

bench:
	@go test -bench=. -benchmem -run=NONE ./...

prof:
ifndef PKG 
	@echo 'missing environment variables: PKG'
	@exit 1
else
	@if [ ! -d "./internal/${PKG}" ]; then\
		echo 'not found: PKG';\
		exit 1;\
	fi

	@for t in cpu block mem; do\
		go test -bench=. -run=NONE -$${t}profile=$${t}.out ./internal/${PKG};\
		go tool pprof -text -nodecount=10 ./${PKG}.test $${t}.out;\
	done
endif

lint:
	@go vet ./...
	@golint -set_exit_status ./...

clean:
	@rm -f ${APP_NAME} main

run:
	@BRAIN_URL=http://127.0.0.1:${BRAIN_PORT} ./${APP_NAME}

run-debug:
	@BRAIN_URL=http://127.0.0.1:${BRAIN_PORT} ./${APP_NAME} -debug

build-image:
	@docker build -t ${APP_NAME} .
	@docker image prune -f

lint-image:
	@docker run --rm -i hadolint/hadolint < Dockerfile

run-container:
	@docker run --env-file=.env --rm ${APP_NAME}

clean-image:
	@docker rmi -f ${APP_NAME}

run-brain:
	@docker-compose -f docker-compose.development.yml up

stop-brain:
	@docker-compose -f docker-compose.development.yml down

redis-aof:
	@docker-compose exec redis redis-cli bgrewriteaof

.PHONY: all build test lint clean run build-image lint-image run-container clean-image run-brain stop-brain
