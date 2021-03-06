# Elasticsearch-Batcher

[![Build Status](https://travis-ci.org/shujew/elasticsearch-batcher.svg?branch=master)](https://travis-ci.org/shujew/elasticsearch-batcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/shujew/elasticsearch-batcher)](https://goreportcard.com/report/github.com/shujew/elasticsearch-batcher)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Tested using golang v1.14.0

## Introduction
Elasticsearch-Batcher is a a FIFO queue for indexing documents to Elasticsearch. It uses batch processing and guarantees
ordering at the client-level*. That means that when deployed across several instances, all of client's A request will be
processed in order but there is no guarantee that it will be processed before client B's requests, which may have arrived
after client A's request.\
\
\* this assumes that you can route requests from a client to the same container (e.g. if you are using AWS, you can set
up your target group and load balancer to do so - search for stickiness)

## Why I built Elasticsearch-Batcher
I was initially toying around with AWS Firehose Delivery Streams to Elasticsearch but I hit a wall when I wanted to specify
the `_id` field of a document when indexing it. Doing further research, it turned out that AWS Firehose Delivery Streams did
not support that feature and were missing other features as well (such as the update bulk operation). Thus, I wrote
Elasticsearch-Batcher, aimed to be a proxy between a client and an ES server, allowing for batch processing with the full
potential of the `_bulk` endpoint, which AWS Firehose Delivery Streams do not support.

## Getting Started

### Sending Data

Simply send the same data you would normally send to the [elasticsearch bulk api](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html)
to `/ingest/v1`

### Configuration 

Elasticsearch-Batcher is easily configurable via global variables:

- `ESB_DEBUG`
  - `true` to enable verbose logging 
  - Defaults to `false`
- `ESB_HTTP_PORT`
  - Set it to desired port you want app to run on
  - Defaults to `8889`
- `ESB_ALLOW_ALL_ORIGINS`
  - Set to `true` to allow any origin (CORS) 
  - Defaults to `true`
- `ESB_ALLOWED_ORIGINS`
  - Comma separated list of allowed origins (CORS) if `ESB_ALLOW_ALL_HOSTS=false`
  - Defaults to an empty string
- `ESB_ES_HOST`
  - protocol + hostname of ES cluster (e.g. http://mycluster.com)
  - Defaults to `http://localhost:9200`
- `ESB_ES_USERNAME`
  - es cluster username if any (for basic auth)
  - Defaults to an empty string
- `ESB_ES_PASSWORD`
  - es cluster password if any (for basic auth)
  - Defaults to an empty string
- `ESB_ES_TIMEOUT_SECONDS`
  - Set to how long you wish to give ES to ingest data (in seconds)
  - Defaults to `60`
- `ESB_FLUSH_INTERVAL_SECONDS`
  - Set to desired value after which events should be flushed to ES (in seconds)
  - Defaults to `60`
  
### Running (locally)

#### Pre-requisites
- [golang](https://golang.org/doc/install)
- [dep](https://github.com/golang/dep)

#### Extraction
Ensure files are extracted to `$GOPATH/src/github.com/shujew/elasticsearch-batcher/`

#### Installing
```shell script
make install
```

#### Start server
For this to work, $GOBIN must be included in your $PATH
```shell script
elasticsearch-batcher
```

#### Development
Running this way is meant for development purpose only. It shouldn't be used in production as it compiles the latest
version then runs which can cause delays when starting a server.
```shell script
make run
```

### Running (using docker-compose)

This is intended for easily set up an environment running
- Elasticsearch (http://localhost:9200)
- Kibana (http://localhost:5601)
- Elasticsearch-Batcher (http://localhost:8889)

#### Pre-requisites
- [docker](https://www.docker.com/)

#### Start containers
```shell script
cd /path/to/repo
make run-docker
```

### Helpful docker commands

#### Rebuilding docker container (after a git pull for example)
```shell script
cd /path/to/repo
make image
```

#### Connecting to container
- run `docker ps`
- Grab the container id of desired image from the `CONTAINER ID` column. You can recognize containers by the `NAMES` column (e.g. `elasticsearch:7.6.1`)
- In a terminal window, run `docker exec -it <CONTAINER-ID> /bin/bash`

#### Listing running containers
```shell script
docker ps
```

#### Stopping all containers
```shell script
docker stop $(docker ps -aq)
```

#### Removing all containers
```shell script
docker rm $(docker ps -aq)
```

#### Removing all images
```shell script
docker rmi $(docker images -q)
