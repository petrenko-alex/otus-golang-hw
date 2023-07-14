package hw05parallelexecution

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

type MockTask struct {
	mock.Mock

	TaskDuration time.Duration
}

func (m *MockTask) exec() error {
	time.Sleep(m.TaskDuration)

	args := m.Called()
	return args.Error(0)
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := generateFailedTasks(tasksCount)
		workersCount := 10
		maxErrorsCount := 23

		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(
			t,
			getFinishedMockTaskCount(tasks),
			int32(workersCount+maxErrorsCount),
			"extra tasks were started",
		)
	})

	t.Run("tasks without errors", func(t *testing.T) {
		taskCount := 50
		tasks := generateSuccessTasks(taskCount)
		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, taskCount, getFinishedMockTaskCount(tasks), "not all tasks were completed")
		require.LessOrEqual(
			t,
			int64(elapsedTime),
			int64(getMockTasksDuration(tasks)/2),
			"tasks were run sequentially?",
		)

		// TODO: fix time suggestion assert
	})

	t.Run("process all tasks, have some errors", func(t *testing.T) {
		// Повседневный сценарий: есть какие-то ошибки, но лимит не превышен, все задачи выполнены

	})

	t.Run("no more than N goroutines run", func(t *testing.T) {
		// Проверка, что запускается не более N горутин
	})

	t.Run("stop on M errors", func(t *testing.T) {
		// Проверка, что выполнится ровно M задач при 1 воркере
	})

	t.Run("all goroutines stopped", func(t *testing.T) {
		// Проверка, что не осталось запущенных горутин
	})

	t.Run("all tasks have same exec time", func(t *testing.T) {
		// Тест на случай, когда все таски выполняются одинаковое время
	})

	t.Run("tasks count lass than workers count", func(t *testing.T) {
		// Количество задач, меньше количества воркеров
	})

	t.Run("tasks count equals workers count", func(t *testing.T) {
		// Количество задач равно количеству воркеров
	})

	t.Run("No errors allowed (M equals zero)", func(t *testing.T) {

	})

	t.Run("Errors ignored (M less than zero)", func(t *testing.T) {

	})
}

func generateSuccessTasks(tasksCount int) []ExecutableTask {
	tasks := make([]ExecutableTask, 0, tasksCount)

	for i := 0; i < tasksCount; i++ {
		task := &MockTask{
			TaskDuration: time.Millisecond * time.Duration(rand.Intn(100)),
		}
		task.On("exec").Return(nil)
		tasks = append(tasks)
	}

	return tasks
}

func generateFailedTasks(tasksCount int) []ExecutableTask {
	tasks := make([]ExecutableTask, 0, tasksCount)

	for i := 0; i < tasksCount; i++ {
		task := &MockTask{
			TaskDuration: time.Millisecond * time.Duration(rand.Intn(100)),
		}

		task.On("exec").Return(fmt.Errorf("error from task %d", i))
	}

	return tasks
}

func getFinishedMockTaskCount(tasks []ExecutableTask) int {
	var taskCount int

	for _, task := range tasks {
		mockTask, ok := task.(*MockTask)
		if !ok {
			panic("tasks should be MockTask")
		}

		taskCount += len(mockTask.Calls)
	}

	return taskCount
}

func getMockTasksDuration(tasks []ExecutableTask) time.Duration {
	var tasksDuration time.Duration

	for _, task := range tasks {
		mockTask, ok := task.(*MockTask)
		if !ok {
			panic("tasks should be MockTask")
		}

		tasksDuration += mockTask.TaskDuration
	}

	return tasksDuration
}
