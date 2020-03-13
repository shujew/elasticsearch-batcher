package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"shujew/elasticsearch-batcher/batch"
	"time"
)

type BulkClient struct {
	esHost        string
	esUsername    string
	esPassword    string

	httpClient    *http.Client
	memoryBatcher *batch.MemoryBatcher
}

var clientSingleton = newBulkClient(
	"localhost",
	30 * time.Second,
	10 * time.Second,
)

func newBulkClient(
	esHost        string,
	httpTimeout   time.Duration,
	flushInterval time.Duration,
	) *BulkClient {

	log.WithFields(log.Fields{
		"esHost": esHost,
		"timeout": httpTimeout,
		"flushInterval": flushInterval,
	}).Debug("creating new bulk es client")

	client := BulkClient{
		esHost: esHost,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
		memoryBatcher: batch.NewMemoryBatch(flushInterval),
	}

	client.Start()

	return &client
}

func GetBulkClient() *BulkClient {
	return clientSingleton
}

func (c *BulkClient) SetBasicAuth(esUsername, esPassword string) {
	log.WithFields(log.Fields{
		"username": esUsername,
		"password": esPassword,
	}).Debug("setting basic auth on es client")

	c.esUsername = esUsername
	c.esPassword = esPassword
}

func (c *BulkClient) Start() {
	log.Debug("starting bulk es client")

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

func (c *BulkClient) Stop() {
	log.Debug("stopping bulk es client")

	c.memoryBatcher.Stop()
}

func (c *BulkClient) IndexDocument(document map[string]interface{}) {
	log.Debug("sending item to be indexed to bulk es client")
	// actually schedules document for bulk indexing
	c.memoryBatcher.AddItem(document)
}

func (c *BulkClient) bulkIndexDocuments(documents []interface{}) {
	log.Debug("indexing documents from bulk es client to es host")

	endpoint := fmt.Sprintf("%s/_bulk", c.esHost)
	body := c.generateBulkPayload(documents)

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if c.esUsername != "" && c.esPassword != "" {
		req.SetBasicAuth(c.esUsername, c.esPassword)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")

	// TODO: add support for logging failures or atleast retries
	go c.httpClient.Do(req)
}

func (c *BulkClient) generateBulkPayload(documents []interface{}) []byte {
	payload := ""

	for _, document := range documents {
		if documentPayload, err := json.Marshal(document); err == nil {
			// leaving adding newline to the client
			payload = payload + string(documentPayload)
		}
	}

	return []byte(payload)
}

