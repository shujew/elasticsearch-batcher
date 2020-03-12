package batch

import (
	"sync"
	"time"
)

type MemoryBatcher struct {
	items    []interface{}
	interval time.Duration

	mutex sync.RWMutex

	JobsChan chan []interface{}
	quitChan chan struct{}
}

func NewMemoryBatch(interval time.Duration) *MemoryBatcher {
	return &MemoryBatcher{
		items:    []interface{}{},
		interval: interval,
	}
}

func (b *MemoryBatcher) Start() {
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

func (b *MemoryBatcher) Stop() {
	b.flushItems(func() {
		defer close(b.JobsChan)
		defer close(b.quitChan)

		b.quitChan <- struct{}{}
	})
}

func (b *MemoryBatcher) AddItem(item interface{}) {
	b.mutex.Lock()
	go func(item interface{}) {
		defer b.mutex.Unlock()
		b.items = append(b.items, item)
	}(item)
}

func (b *MemoryBatcher) flushItems(completion func()) {
	b.mutex.Lock()
	go func(items []interface{}) {
		defer b.mutex.Unlock()
		if len(b.items) > 0 {
			tmp := make([]interface{}, len(b.items))
			copy(tmp, b.items)
			b.JobsChan <- tmp
			b.items = []interface{}{}

			if completion != nil {
				completion()
			}
		}
	}(b.items)
}
