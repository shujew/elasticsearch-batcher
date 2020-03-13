# Start from a Debian image with the Go 1.14.x installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.14

# Copy the local package files to the container's workspace.
ADD . /go/src/shujew/elasticsearch-batcher/

RUN \
    # install dep (dependency management)
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh \
    # cd into elasticsearch-batcher directory
    && cd /go/src/shujew/elasticsearch-batcher/ \
    # install dependencies
    && dep ensure

# Build elastic-batcher
RUN go install shujew/elasticsearch-batcher
