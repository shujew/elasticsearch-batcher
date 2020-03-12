package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func NewBulkClient(
	esHost        string,
	httpTimeout   time.Duration,
	flushInterval time.Duration,
	) *BulkClient {

	client := BulkClient{
		esHost: esHost,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
		memoryBatcher: batch.NewMemoryBatch(flushInterval),
	}

	return &client
}

func (c *BulkClient) SetBasicAuth(esUsername, esPassword string) {
	c.esUsername = esUsername
	c.esPassword = esPassword
}

func (c *BulkClient) Start() {
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
	c.memoryBatcher.Stop()
}

func (c *BulkClient) IndexDocument(document map[string]interface{}) {
	// actually schedules document for bulk indexing
	c.memoryBatcher.AddItem(document)
}

func (c *BulkClient) bulkIndexDocuments(documents []interface{}) {
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

