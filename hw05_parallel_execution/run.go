package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type ExecutableTask interface {
	Exec() error
}

type Task struct {
	task func() error
}

func (t Task) Exec() error {
	return t.task()
}

type TaskRunner struct {
	maxErrors     int
	workersCount  int
	errorsCounter int

	tasks       []ExecutableTask
	taskChannel chan ExecutableTask

	sync.WaitGroup
	sync.Mutex
}

func NewTaskRunner(tasks []ExecutableTask, workersCount, maxErrors int) *TaskRunner {
	return &TaskRunner{
		maxErrors:    maxErrors,
		workersCount: workersCount,
		tasks:        tasks,

		taskChannel: make(chan ExecutableTask, len(tasks)),
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func (r *TaskRunner) Run() error {
	r.Add(r.workersCount)

	r.runWorkers()
	r.sendTasks()

	close(r.taskChannel)
	r.Wait()

	if r.errorLimitExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func (r *TaskRunner) getErrorsCounter() int {
	r.Lock()
	defer r.Unlock()

	return r.errorsCounter
}

func (r *TaskRunner) incErrorsCounter() {
	r.Lock()
	r.errorsCounter++
	r.Unlock()
}

func (r *TaskRunner) runWorkers() {
	for i := 0; i < r.workersCount; i++ {
		go func() { r.runWorker() }()
	}
}

func (r *TaskRunner) runWorker() {
	defer r.Done()

	for {
		if r.maxErrors >= 0 && r.getErrorsCounter() >= r.maxErrors {
			break
		}

		task, continueWork := <-r.taskChannel
		if !continueWork {
			break
		}

		taskError := task.Exec()
		if taskError != nil {
			r.incErrorsCounter()
		}
	}
}

func (r *TaskRunner) sendTasks() {
	for _, task := range r.tasks {
		r.taskChannel <- task
	}
}

func (r *TaskRunner) errorLimitExceeded() bool {
	return r.maxErrors >= 0 && r.getErrorsCounter() >= r.maxErrors
}
