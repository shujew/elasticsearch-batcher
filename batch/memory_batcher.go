// Package batch provides structs to perform batch processing
package batch

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// MemoryBatcher encapsulates channels and mutexes to provide a
// clean and easy-to-use FIFO queue to batch data in memory and
// emit them at a specified interval
type MemoryBatcher struct {
	items    []interface{}
	interval time.Duration

	mutex sync.RWMutex

	JobsChan chan []interface{}
	quitChan chan struct{}
}

// NewMemoryBatcher creates new instance of instance of
// MemoryBatcher configured to emit events a specific
// interval, and returns a pointer to it
func NewMemoryBatcher(interval time.Duration) *MemoryBatcher {
	log.WithFields(log.Fields{
		"interval": interval,
	}).Trace("creating new memory batcher")

	return &MemoryBatcher{
		items:    []interface{}{},
		interval: interval,
	}
}

// Start sets up the channels required to queue data
// and terminate the instance. It also sets up data
// flushing after the specified interval
func (b *MemoryBatcher) Start() {
	log.Trace("starting memory batcher")

	b.JobsChan = make(chan []interface{})
	b.quitChan = make(chan struct{})

	ticker := time.NewTicker(b.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				b.flushItems(nil)
			case <-b.quitChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop closes all all open channels and cleans
// up used resources
func (b *MemoryBatcher) Stop() {
	log.Trace("stopping memory batcher")

	b.flushItems(func() {
		defer close(b.JobsChan)
		defer close(b.quitChan)

		b.quitChan <- struct{}{}
	})
}

// AddItem is an atomic function which adds a
// type-agnostic item to the queue
func (b *MemoryBatcher) AddItem(item interface{}) {
	log.Trace("adding item to memory batcher")

	b.mutex.Lock()
	go func(item interface{}) {
		defer b.mutex.Unlock()
		b.items = append(b.items, item)
	}(item)
}

// flushItems is an atomic function which handles
// flushing the queue to jobs channel (and runs
// a callback afterwards if specified)
func (b *MemoryBatcher) flushItems(completion func()) {
	b.mutex.Lock()
	go func(items []interface{}) {
		defer b.mutex.Unlock()
		if len(b.items) > 0 {
			log.WithFields(log.Fields{
				"count": len(b.items),
			}).Debug("flushing items from memory batcher")

			tmp := make([]interface{}, len(b.items))
			copy(tmp, b.items)
			b.JobsChan <- tmp
			b.items = []interface{}{}

			if completion != nil {
				log.Trace("calling completion after flushing items")
				completion()
			}
		}
	}(b.items)
}
