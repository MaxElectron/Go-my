//go:build !solution

package rwmutex

type RWMutex struct {
	rlock chan int
	wlock chan bool
}

func New() *RWMutex {
	return &RWMutex{
		rlock: make(chan int, 1),
		wlock: make(chan bool, 1),
	}
}

func (rw *RWMutex) RLock() {
	select {
	case rw.wlock <- false:
		rw.rlock <- 1
	case readerCount := <-rw.rlock:
		rw.rlock <- readerCount + 1
	}
}

func (rw *RWMutex) RUnlock() {
	readerCount := <-rw.rlock
	readerCount--

	if readerCount != 0 {
		rw.rlock <- readerCount
	} else {
		<-rw.wlock
	}
}

func (rw *RWMutex) Lock() {
	rw.wlock <- true
}

func (rw *RWMutex) Unlock() {
	<-rw.wlock
}
