package v1

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"shujew/elasticsearch-batcher/data/elasticsearch"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	setDefaultHeaders(&w, req)

	switch req.Method {
	case "POST":
		POSTHandler(w, req)
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func POSTHandler(w http.ResponseWriter, req *http.Request) {
	if body, err := ioutil.ReadAll(req.Body); err == nil {
		esClient := elasticsearch.GetBulkClient()
		esClient.QueueForBulkIndexing(body)
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
