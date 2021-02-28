package pkg

// Worker holds executes tasks in jobQueue
type Worker struct {
	jobChannel  chan interface{}
	closeHandle chan bool
}

// NewWorker creates the new worker
func NewWorker(jobChannel chan interface{}, closeHandle chan bool) *Worker {
	return &Worker{jobChannel: jobChannel, closeHandle: closeHandle}
}

// Start starts the worker by listening to the job channel
func (w Worker) Start() {
	go func() {
		for {
			select {
			// Wait for the job
			case job := <-w.jobChannel:
				// Got the job
				w.executeJob(job)
			case <-w.closeHandle:
				// Exit the go routine when the closeHandle channel is closed
				return
			}
		}
	}()
}

// executeJob executes the job based on the type
func (w Worker) executeJob(job interface{}) {
	// Execute the job based on the task type
	switch task := job.(type) {
	case RunnableTask:
		task.Task.Run()
		break
	case CallableTask:
		response := task.Task.Call()
		task.Handle.done = true
		task.Handle.response <- response
		break
	}
}
