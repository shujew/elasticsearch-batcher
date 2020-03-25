// Package ingestv1 provides an implementation of the
// /ingest/v1 endpoint to queue data for indexing
// to an elasticsearch cluster
package ingestv1

import (
	"encoding/json"
	"github.com/shujew/elasticsearch-batcher/elasticsearch"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// Handler sets the default headers for CORS on the
// request and routes it to its proper handle based
// on the method used
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

// POSTHandler queues up the body of the request for
// later indexing into an  elasticsearch cluster
func POSTHandler(w http.ResponseWriter, req *http.Request) {
	if body, err := ioutil.ReadAll(req.Body); err == nil {
		if len(body) == 0 {
			//TODO: add a sync endpoint instead to synchronize timestamps
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(generateSyncResponse())
		} else {
			esClient := elasticsearch.GetBulkClient()
			esClient.QueueForBulkIndexing(body)
			w.WriteHeader(http.StatusCreated)
		}
	} else {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func generateSyncResponse() []byte {
	response := map[string]int64{
		"t": time.Now().Unix(),
	}

	if responseString, err := json.Marshal(response); err == nil {
		return responseString
	}

	return []byte{}
}
