package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

//TODO: allow configure these via env vars
var httpPort = 8889
var logLevel = log.DebugLevel

func main() {
	configureLogging()

	httpAddr := fmt.Sprintf(":%d", httpPort)
	log.Info("server is listening on port ", httpPort)

	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		panic(err)
	}
}

func configureLogging() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	Formatter.ForceColors = true
	log.SetFormatter(Formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)
}