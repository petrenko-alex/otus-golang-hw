package hw05parallelexecution

import (
	"errors"
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
func Run(tasks []ExecutableTask, n, m int) error {
	// Place your code here.

	return nil
}
