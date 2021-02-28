package workerpool

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkerpool(t *testing.T) {
	wp := NewWorkerpool(1, 1)
	assert.NotNil(t, wp)
	wp.Close()
	_, ok := <-wp.jobQueue
	assert.False(t, ok)
	_, ok = <-wp.closeHandle
	assert.False(t, ok)
}

type taskRunnable struct {
	wg    *sync.WaitGroup
	value string
}

func (t *taskRunnable) Run() {
	t.value = "pass"
	t.wg.Done()
}

func TestWorkerpoolExecute(t *testing.T) {
	wp := NewWorkerpool(1, 1)
	defer wp.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	task := &taskRunnable{wg: &wg, value: "fail"}
	assert.NotEqual(t, "pass", task.value)

	wp.Execute(task)
	wg.Wait()
	assert.Equal(t, "pass", task.value)
}

type taskCallable struct{}

func (t *taskCallable) Call() interface{} {
	return "pass"
}

func TestWorkerpoolExecuteFuture(t *testing.T) {
	wp := NewWorkerpool(1, 1)
	defer wp.Close()

	task := &taskCallable{}
	f, _ := wp.ExecuteFuture(task)
	assert.Equal(t, "pass", f.Get())
	assert.True(t, f.IsDone())
}

func TestWorkerpoolExecuteQueueFullErr(t *testing.T) {
	wp := NewWorkerpool(0, 1)
	defer wp.Close()
	task := &taskRunnable{}
	err := wp.Execute(task)
	assert.Nil(t, err)

	err = wp.Execute(task)
	assert.Equal(t, ErrQueueFull, err)
}

func TestWorkerpoolExecuteFutureQueueFullErr(t *testing.T) {
	wp := NewWorkerpool(0, 1)
	defer wp.Close()
	task := &taskCallable{}
	_, err := wp.ExecuteFuture(task)
	assert.Nil(t, err)

	_, err = wp.ExecuteFuture(task)
	assert.Equal(t, ErrQueueFull, err)
}
