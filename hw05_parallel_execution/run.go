package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type ExecutableTask interface {
	exec() error
}

type Task struct {
	task func() error
}

func (t Task) exec() error {
	return t.task()
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []ExecutableTask, workersCount, maxErrors int) error {
	taskChannel := make(chan ExecutableTask) // todo: need buffer ?
	wg := sync.WaitGroup{}
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go func() {
			defer wg.Done()

			for {
				task, continueWork := <-taskChannel
				if !continueWork {
					break
				}

				taskError := task.exec()
				if taskError != nil {
					// todo:
				}
			}
		}()
	}

	for _, task := range tasks {
		taskChannel <- task
	}
	close(taskChannel)

	wg.Wait()

	return nil
}
