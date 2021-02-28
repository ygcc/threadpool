package pkg

// Callable is the task which has returns
type Callable interface {
	Call() interface{}
}

// Future is the handle returned after enqueuing a callable task
type Future struct {
	response chan interface{}
	done     bool
}

// NewFuture returns new future
func NewFuture() *Future {
	return &Future{response: make(chan interface{})}
}

// CallableTask wraps the callable and future together
type CallableTask struct {
	Task   Callable
	Handle *Future
}

// Get returns the response of the Callable task when done
// It is blocking call and waits for the execution to complete
func (f *Future) Get() interface{} {
	return <-f.response
}

// IsDone returns true if the execution is already done
func (f *Future) IsDone() bool {
	return f.done
}
