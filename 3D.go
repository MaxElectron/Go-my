//go:build !solution

package waitgroup

// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.

type WaitCounter struct {
	counter int
	wait    chan struct{}
}

type WaitGroup struct {
	waitCount chan WaitCounter
}

func MakeWaitCounter() *WaitCounter {
	return &WaitCounter{
		counter: 0,
		wait:    make(chan struct{}),
	}
}

// New creates WaitGroup.
func New() *WaitGroup {
	waitGroup := &WaitGroup{
		waitCount: make(chan WaitCounter, 1),
	}
	newWaitCounter := MakeWaitCounter()
	close(newWaitCounter.wait)
	waitGroup.waitCount <- *newWaitCounter
	return waitGroup
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
	waitCounter := <-wg.waitCount

	if waitCounter.counter == 0 {
		waitCounter = *MakeWaitCounter()
	}

	waitCounter.counter += delta

	if waitCounter.counter < 0 {
		panic("negative WaitGroup counter")
	}

	if waitCounter.counter == 0 {
		close(waitCounter.wait)
	}

	wg.waitCount <- waitCounter
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	waitCounter := <-wg.waitCount
	wait := waitCounter.wait
	wg.waitCount <- waitCounter

	<-wait
}
