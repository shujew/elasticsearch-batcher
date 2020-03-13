# Elasticsearch-Batcher

## Intro
Elasticsearch-Batcher is a a fifo queue for Elasticsearch `_bulk` endpoint, which
guarantees ordering at the client-level.

## Sending Data

Endpoint: `/ingest/v1` \
Data: Use the same format as [Elasticsearch bulk api](https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html)

## Configuration

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

## Local Setup

### Pre-requisites
- [Docker](https://www.docker.com/)

### Run Elastic-Batcher
```shell script
cd /path/to/repo
docker-compose up --build
```

Elasticsearch-Batcher will accessible at http://localhost:8889/ \
Elasticsearch will accessible at http://localhost:9200/ \
Kibana will accessible at http://localhost:5601/

### Helpful docker commands

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