package threadpool

import (
	"sync"
	"time"
	"github.com/shettyh/threadpool/internal"
)

// ScheduledThreadPool
// Schedules the task with the given delay
type ScheduledThreadPool struct {
	workers     chan chan interface{}
	tasks       *sync.Map
	noOfWorkers int
	counter     uint64
	counterLock sync.Mutex
}

// NewScheduledThreadPool creates new scheduler thread pool with given number of workers
func NewScheduledThreadPool(noOfWorkers int) *ScheduledThreadPool {
	pool := &ScheduledThreadPool{}
	pool.noOfWorkers = noOfWorkers
	pool.workers = make(chan chan interface{}, noOfWorkers)
	pool.tasks = new(sync.Map)
	pool.createPool()
	return pool
}

// createPool creates the workers pool
func (stf *ScheduledThreadPool) createPool() {
	for i := 0; i < stf.noOfWorkers; i++ {
		worker := NewWorker(stf.workers)
		worker.Start()
	}

	go stf.dispatch()
}

// dispatch will check for the task to run for current time and invoke the task
func (stf *ScheduledThreadPool) dispatch() {
	for {
		go stf.intervalRunner()     // Runner to check the task to run for current time
		time.Sleep(time.Second * 1) // Check again after 1 sec
	}
}

// intervalRunner checks the tasks map and runs the tasks that are applicable at this point of time
func (stf *ScheduledThreadPool) intervalRunner() {
	// update the time count
	stf.updateCounter()

	// Get the task for the counter value
	currentTasksToRun, ok := stf.tasks.Load(stf.counter)

	// Found tasks
	if ok {
		// Convert to tasks set
		currentTasksSet := currentTasksToRun.(*internal.Set)

		// For each tasks , get a worker from the pool and run the task
		for _, val := range currentTasksSet.GetAll() {
			go func(job interface{}) {
				// get the worker from pool who is free
				worker := <-stf.workers
				// Submit the job to the worker
				worker <- job
			}(val)
		}
	}
}

// updateCounter thread safe update of counter
func (stf *ScheduledThreadPool) updateCounter() {
	stf.counterLock.Lock()
	stf.counter++
	stf.counterLock.Unlock()
}

// ScheduleOnce the task with given delay
func (stf *ScheduledThreadPool) ScheduleOnce(task Runnable, delay time.Duration) {
	scheduleTime := stf.counter + uint64(delay.Seconds())
	existingTasks, ok := stf.tasks.Load(scheduleTime)

	// Create new set if no tasks are already there
	if !ok {
		existingTasks = internal.NewSet()
		stf.tasks.Store(scheduleTime, existingTasks)
	}
	// Add task
	existingTasks.(*internal.Set).Add(task)
}