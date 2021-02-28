package workerpool

import (
	"fmt"

	pkg "github.com/ygcc/workerpool/pkg"
)

var (
	// ErrQueueFull indicates task queue is full
	ErrQueueFull = fmt.Errorf("Task queue is full, can not add task")
)

// Workerpool receives jobs and sends to workers
type Workerpool struct {
	queueSize   int64
	noOfWorkers int

	jobQueue    chan interface{}
	closeHandle chan bool // Channel used to stop all the workers
}

// NewWorkerpool creates workerpool
func NewWorkerpool(noOfWorkers int, queueSize int64) *Workerpool {
	threadPool := &Workerpool{queueSize: queueSize, noOfWorkers: noOfWorkers}
	threadPool.jobQueue = make(chan interface{}, queueSize)
	threadPool.closeHandle = make(chan bool)
	threadPool.createPool()
	return threadPool
}

func (t *Workerpool) submitTask(task interface{}) error {
	select {
	case t.jobQueue <- task:
		return nil
	default:
		return ErrQueueFull
	}
}

// Execute enqueues runnable tasks
func (t *Workerpool) Execute(task pkg.Runnable) error {
	return t.submitTask(pkg.RunnableTask{Task: task})
}

// ExecuteFuture enqueues callable tasks and returns the response handle
func (t *Workerpool) ExecuteFuture(task pkg.Callable) (*pkg.Future, error) {
	futureTask := pkg.CallableTask{Task: task, Handle: pkg.NewFuture()}
	err := t.submitTask(futureTask)
	if err != nil {
		return nil, err
	}
	return futureTask.Handle, nil
}

// Close will close the threadpool
// It sends the stop signal to all the worker that are running
func (t *Workerpool) Close() {
	close(t.closeHandle) // Stops all the routines
	close(t.jobQueue)    // Closes the job Queue
}

// createPool creates the workers which start listening on the jobQueue
func (t *Workerpool) createPool() {
	for i := 0; i < t.noOfWorkers; i++ {
		worker := pkg.NewWorker(t.jobQueue, t.closeHandle)
		worker.Start()
	}
}
