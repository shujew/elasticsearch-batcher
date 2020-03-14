package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/shujew/elasticsearch-batcher/config"
	"github.com/shujew/elasticsearch-batcher/routing/ingestv1"
)

func main() {
	// setting up logging
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	logLevel := config.GetLogLevel()
	log.SetLevel(logLevel)

	// setting up routes
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte("Elasticsearch Batcher"))
	})
	http.HandleFunc("/ingest/v1", ingestv1.Handler)

	// setting up http server
	httpPort := config.GetHTTPPort()
	httpAddr := fmt.Sprintf(":%s", httpPort)
	log.Info("server is listening on port ", httpPort)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		panic(err)
	}
}
