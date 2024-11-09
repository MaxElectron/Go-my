//go:build !solution

package batcher

import (
	"sync"

	"gitlab.com/slon/shad-go/batcher/slow"
)

type Batcher struct {
	mu             *sync.Mutex
	batcherCond    *sync.Cond
	valueReadyCond *sync.Cond
	value          *slow.Value
	inter          interface{}
	isRunning      bool
}

func (batcher *Batcher) Load() interface{} {
	batcher.mu.Lock()
	defer batcher.mu.Unlock()

	for {
		if !batcher.isRunning {
			batcher.isRunning = true
			batcher.batcherCond.Broadcast()

			go batcher.loadValue()

			batcher.valueReadyCond.Wait()
			batcher.isRunning = false
			batcher.batcherCond.Broadcast()
		} else {
			batcher.batcherCond.Wait()
			if !batcher.isRunning {
				continue
			}
			batcher.valueReadyCond.Wait()
		}

		return batcher.inter
	}
}

func (batcher *Batcher) loadValue() {
	batcher.inter = batcher.value.Load()
	batcher.valueReadyCond.Broadcast()
}

func NewBatcher(value *slow.Value) *Batcher {
	mu := &sync.Mutex{}
	return &Batcher{
		mu:             mu,
		batcherCond:    sync.NewCond(mu),
		valueReadyCond: sync.NewCond(mu),
		value:          value,
		isRunning:      false,
	}
}
