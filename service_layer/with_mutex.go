package servicelayer

import (
	"sync"
	"time"
)

type DummyModel struct {
	ID    int64
	Stock int64
	Mu    sync.RWMutex // that's the key
}

func (d *DummyModel) Writer(fakeQuantity int64) int64 {
	// implementation
	/*
		here we use Read-Write Mutex so can block if another process in memory try to write it, until this process done
	*/

	// [*] add scenario latency
	time.Sleep(500 * time.Millisecond)

	d.Mu.Lock()
	calcualted := d.Stock - fakeQuantity
	d.Stock = calcualted
	d.Mu.Unlock()
	return calcualted
}

func (d *DummyModel) Reader() int64 {
	d.Mu.RLock()
	Stock := &d.Stock
	d.Mu.RUnlock()
	return *Stock
}
