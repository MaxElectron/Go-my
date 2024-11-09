//go:build !solution

package dupcall

import (
	"context"
	"sync"
	"sync/atomic"
)

type WaitGroupCount struct {
	sync.WaitGroup

	cnt int64
}

func (wg *WaitGroupCount) Done() {
	atomic.AddInt64(&wg.cnt, -1)
	wg.WaitGroup.Done()
}

func (wg *WaitGroupCount) Load() int {
	return int(atomic.LoadInt64(&wg.cnt))
}

type Call struct {
	mu       sync.Mutex
	prev_ctx context.Context
	once     sync.Once
	wg       WaitGroupCount
	ctx      context.Context
	cancel   context.CancelFunc

	results chan struct {
		res interface{}
		err error
	}
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (interface{}, error) {

	o.once.Do(func() {
		o.prev_ctx = context.Background()
		o.wg = WaitGroupCount{}
		o.results = make(chan struct {
			res interface{}
			err error
		})
	})

	atomic.AddInt64(&o.wg.cnt, int64(1))
	o.wg.WaitGroup.Add(1)

	if o.mu.TryLock() {
		o.ctx, o.cancel = context.WithCancel(o.prev_ctx)

		go func() {
			res, err := cb(o.ctx)
			for i := 0; i < o.wg.Load(); i++ {
				o.results <- struct {
					res interface{}
					err error
				}{res, err}
			}
			o.mu.Unlock()
		}()
	}

	select {
	case <-ctx.Done():
		o.wg.Done()
		o.cancel()
		return nil, ctx.Err()

	case res := <-o.results:
		o.wg.Done()
		return res.res, res.err
	}
}
