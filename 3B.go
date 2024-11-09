package once

type Once struct {
	done chan bool
}

func New() *Once {
	o := &Once{
		done: make(chan bool, 1),
	}
	o.done <- false
	return o
}

func (o *Once) Do(f func()) {
	defer func() { o.done <- true }()
	if <-o.done {
		return
	}
	f()
}
