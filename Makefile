# Makefile; used to make build step easier.
# Usage: run `make [task_name]`, e.g. `make install`.
# If there are any issues when updating this file, make sure tabs are used instead of spaces.

build:
	dep ensure && go build github.com/shujew/elasticsearch-batcher

install:
	dep ensure && go install github.com/shujew/elasticsearch-batcher

image:
	docker-compose build

run:
	dep ensure && go run github.com/shujew/elasticsearch-batcher

run-docker:
	docker-compose up