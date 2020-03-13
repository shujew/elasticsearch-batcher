package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	ingestv1 "shujew/elasticsearch-batcher/routing/ingest/v1"
)

//TODO: allow configure these via env vars
var httpPort = 8889
var logLevel = log.TraceLevel

func main() {
	configureLogging()

	httpAddr := fmt.Sprintf(":%d", httpPort)
	log.Info("server is listening on port ", httpPort)

	//Routing
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = w.Write([]byte("Elasticsearch Batcher"))
	})
	http.HandleFunc("/ingest/v1", ingestv1.Handler)

	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		panic(err)
	}
}

func configureLogging() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	log.SetLevel(logLevel)
}