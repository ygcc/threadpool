## Install
`go get github.com/ygcc/workerpool`

## Usage

### Workerpool
- Create instance of `Workerpool` with number of workers required and the size of task queue

  ```
  wp := workerpool.NewWorkerPool(200,1000000)
  ```
  
- Create `Runnable` task and execute
  ```
  type TaskRunnable struct { }
  
  func (t *TaskRunnable) Run(){
    // Do task
  }

  task := &TaskRunnalbe{}
  err := wp.Execute(task)
  ```

- Create `Callable` task and execute
  ```
  type TaskCallable struct { }
  
  func (c *TaskCallable) Call() interface{} {
    //Do task
    return result
  }
  
  //Execute callable task
  task := &TaskCallable{}
  future, err := wp.ExecuteFuture(task)
  
  //Check if the task is done
  isDone := future.IsDone() // true/false
  
  //Get response , blocking call
  result := future.Get()
  ```
- Close the pool
  ```
  wp.Close()
  ```
