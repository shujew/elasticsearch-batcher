// Package elasticsearch provides structs to upload
// documents to an elasticsearch cluster
package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shujew/elasticsearch-batcher/batch"
	"github.com/shujew/elasticsearch-batcher/config"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// BulkClient encapsulates queuing documents
// and emitting them to the _bulk endpoint
// of an elasticsearch cluster
type BulkClient struct {
	esHost     string
	esUsername string
	esPassword string
	esTimeout  time.Duration

	httpClient    *http.Client
	memoryBatcher *batch.MemoryBatcher
}

// clientSingleton is a single instance of BulkClient
// which should be used throughout the whole app
var clientSingleton = newBulkClient(
	config.GetESHost(),
	config.GetESTimeout(),
	config.GetFlushInterval(),
)

// newBulkClient creates and configures a new instance
// of BulkClient
func newBulkClient(
	esHost string,
	httpTimeout time.Duration,
	flushInterval time.Duration,
) *BulkClient {

	log.WithFields(log.Fields{
		"esHost":        esHost,
		"timeout":       httpTimeout,
		"flushInterval": flushInterval,
	}).Trace("creating new bulk es client")

	client := BulkClient{
		esHost: esHost,
		esTimeout: httpTimeout,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
		memoryBatcher: batch.NewMemoryBatcher(flushInterval),
	}

	client.SetBasicAuth(
		config.GetESUsername(),
		config.GetESPassword(),
	)
	client.Start()

	return &client
}

// GetBulkClient returns the singleton instance
// of BulkClient
func GetBulkClient() *BulkClient {
	return clientSingleton
}

// SetBasicAuth sets the username and password used
// when sending data to the elasticsearch cluster
// (BASICAUTH)
func (c *BulkClient) SetBasicAuth(esUsername, esPassword string) {
	log.WithFields(log.Fields{
		"username": esUsername,
		"password": esPassword,
	}).Trace("setting basic auth on es client")

	c.esUsername = esUsername
	c.esPassword = esPassword
}

// Start starts the memory batch and starts a
// goroutine to pick up documents from the
// queue and bulk indexing them
func (c *BulkClient) Start() {
	log.Trace("starting bulk es client")

	c.memoryBatcher.Start()

	go func() {
		for {
			items := <-c.memoryBatcher.JobsChan
			if len(items) == 0 {
				break
			}
			c.bulkIndexDocuments(items)
		}
	}()
}

// Stop terminates the queue / cleans up used resources
func (c *BulkClient) Stop() {
	log.Trace("stopping bulk es client")

	c.memoryBatcher.Stop()
}

// QueueForBulkIndexing adds a document to the queue for
// bulk indexing
func (c *BulkClient) QueueForBulkIndexing(document interface{}) {
	log.Trace("sending item to be indexed to bulk es client")
	// actually schedules document for bulk indexing
	c.memoryBatcher.AddItem(document)
}

// bulkIndexDocuments get a bulk payload and pushes it to the bulk
// endpoint of the elasticsearch cluster
func (c *BulkClient) bulkIndexDocuments(documents []interface{}) {
	log.Trace("indexing documents from bulk es client to es host")

	// sending request to ES
	endpoint := fmt.Sprintf("%s/_bulk?timeout=%ds", c.esHost, c.esTimeout/time.Second)
	reqBody := c.generateBulkPayload(documents)
	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(reqBody))
	if c.esUsername != "" && c.esPassword != "" {
		req.SetBasicAuth(c.esUsername, c.esPassword)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.WithFields(log.Fields{
			"endpoint": endpoint,
			"status_code": resp.StatusCode,
		}).Error("request failed")
		if respBody, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Error(string(respBody))
		}
		return
	}

	// check ES response for errors
	if respBody, err := ioutil.ReadAll(resp.Body); err == nil {
		var respBodyJSON map[string]interface{}
		if err := json.Unmarshal(respBody, &respBodyJSON); err == nil {
			if errPresent, ok := respBodyJSON["errors"].(bool); ok && errPresent {
				log.Error("es host reported errors with payload")
				log.Error(string(reqBody))
				log.Error(string(respBody))
			}
		}
	} else {
		log.Error(err)
	}
}

// generateBulkPayload builds a bulk payload from the emitted documents
// from the queue
func (c *BulkClient) generateBulkPayload(documents []interface{}) []byte {
	var payload []byte

	for _, document := range documents {
		if b, ok := document.([]byte); ok {
			payload = append(payload, b...)
		}
	}

	return payload
}
