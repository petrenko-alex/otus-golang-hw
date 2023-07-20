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
	// todo: need buffer channel ?
	// todo: stop by signal from goroutine
	var errorsCounter int
	mutex := sync.Mutex{}
	taskChannel := make(chan ExecutableTask, len(tasks))

	wg := sync.WaitGroup{}
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go func() {
			defer wg.Done()

			for {
				mutex.Lock()
				eCnt := errorsCounter
				mutex.Unlock()
				if maxErrors >= 0 && eCnt >= maxErrors {
					break
				}

				task, continueWork := <-taskChannel
				if !continueWork {
					break
				}

				taskError := task.exec()
				if taskError != nil {
					mutex.Lock()
					errorsCounter++
					mutex.Unlock()
				}
			}
		}()
	}

	for _, task := range tasks {
		taskChannel <- task
	}
	close(taskChannel)

	wg.Wait()

	if maxErrors >= 0 && errorsCounter >= maxErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}
