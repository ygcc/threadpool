package pkg

// Runnable is the task which has no returns
type Runnable interface {
	Run()
}

// RunnableTask wraps runnable
type RunnableTask struct {
	Task Runnable
}
