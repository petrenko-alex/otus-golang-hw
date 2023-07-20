package hw05parallelexecution

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

type MockTask struct {
	mock.Mock

	TaskDuration      time.Duration
	GoroutinesCounter int
}

func (m *MockTask) exec() error {
	time.Sleep(m.TaskDuration)
	m.GoroutinesCounter = runtime.NumGoroutine()

	args := m.Called()
	return args.Error(0)
}

func NewMockTask() *MockTask {
	return &MockTask{
		TaskDuration: time.Millisecond * time.Duration(rand.Intn(100)),
	}
}

func NewMockTaskWithDuration(duration time.Duration) *MockTask {
	return &MockTask{TaskDuration: duration}
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := generateFailedTasksWithRandomDuration(tasksCount, 100)
		workersCount := 10
		maxErrorsCount := 23

		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(
			t,
			errors.Is(err, ErrErrorsLimitExceeded),
			"actual err - %v", err,
		)
		require.LessOrEqual(
			t,
			getFinishedMockTaskCount(tasks),
			workersCount+maxErrorsCount,
			"extra tasks were started",
		)
	})

	t.Run("tasks without errors", func(t *testing.T) {
		taskCount := 50
		tasks := generateSuccessTasksWithRandomDuration(taskCount)
		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)

		require.NoError(t, err)
		require.Equal(t, taskCount, getFinishedMockTaskCount(tasks), "not all tasks were completed")
		require.LessOrEqual(
			t,
			elapsedTime.Milliseconds(),
			getMockTasksDuration(tasks).Milliseconds()/2,
			"tasks were run sequentially?",
		)
	})

	t.Run("process all tasks, have some errors", func(t *testing.T) {
		taskCount := 10
		workersCount := 4
		maxErrorsCount := 10
		tasks := generateFailedTasksWithRandomDuration(taskCount, 20)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, taskCount, getFinishedMockTaskCount(tasks), "not all tasks were completed")
		assertMockExpectations(t, tasks)
	})

	t.Run("no more than N goroutines run", func(t *testing.T) {
		// Проверка, что запускается не более N горутин
		taskCount := 10
		workersCount := 2
		maxErrorsCount := 50
		tasks := generateSuccessTasksWithRandomDuration(taskCount)

		_ = Run(tasks, workersCount, maxErrorsCount)
		mockTask, ok := tasks[0].(*MockTask)
		if !ok {
			panic("tasks should be MockTask")
		}

		runnerGoroutines := mockTask.GoroutinesCounter - runtime.NumGoroutine()
		require.NotZero(t, mockTask.GoroutinesCounter)
		require.LessOrEqual(t, runnerGoroutines, workersCount)
	})

	t.Run("stop on M errors", func(t *testing.T) {
		// Проверка, что выполнится ровно M задач при 1 воркере
		taskCount := 5
		workersCount := 1
		maxErrorsCount := 2
		tasks := generateFailedTasksWithRandomDuration(taskCount, 100)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.Error(t, err)
		require.Equal(t, maxErrorsCount, getFinishedMockTaskCount(tasks), "not all tasks were completed")
	})

	t.Run("all tasks have same exec time", func(t *testing.T) {
		// Тест на случай, когда все таски выполняются одинаковое время
		taskCount := 10
		workersCount := 4
		maxErrorsCount := 10
		tasks := generateFailedTasks(taskCount, time.Millisecond, 50)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, taskCount, getFinishedMockTaskCount(tasks), "not all tasks were completed")
		assertMockExpectations(t, tasks)
	})

	t.Run("tasks count lass than workers count", func(t *testing.T) {
		// Количество задач, меньше количества воркеров
		taskCount := 1
		workersCount := 10
		maxErrorsCount := 10
		tasks := generateSuccessTasksWithRandomDuration(taskCount)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, taskCount, getFinishedMockTaskCount(tasks), "not all tasks were completed")
		assertMockExpectations(t, tasks)
	})

	t.Run("tasks count equals workers count", func(t *testing.T) {
		// Количество задач равно количеству воркеров
		taskCount := 3
		workersCount := 3
		maxErrorsCount := 10
		tasks := generateSuccessTasksWithRandomDuration(taskCount)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, taskCount, getFinishedMockTaskCount(tasks), "not all tasks were completed")
		assertMockExpectations(t, tasks)
	})

	t.Run("No errors allowed (maxErrors = 0)", func(t *testing.T) {
		taskCount := 10
		workersCount := 4
		maxErrorsCount := 0
		tasks := generateFailedTasksWithRandomDuration(taskCount, 100)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.Error(t, err)
		require.LessOrEqual(t, getFinishedMockTaskCount(tasks), workersCount)
	})

	t.Run("Errors ignored (maxErrors <= 0)", func(t *testing.T) {
		taskCount := 30
		workersCount := 4
		maxErrorsCount := -1
		tasks := generateFailedTasksWithRandomDuration(taskCount, 100)

		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, taskCount, getFinishedMockTaskCount(tasks), "not all tasks were completed")
		assertMockExpectations(t, tasks)
	})

	t.Run("Test concurrency", func(t *testing.T) {
		// TODO: extra task - concurrency w/o time.Sleep
		require.Fail(t, "Not realized")
	})
}

func generateSuccessTasksWithRandomDuration(tasksCount int) []ExecutableTask {
	tasks := make([]ExecutableTask, 0, tasksCount)

	for i := 0; i < tasksCount; i++ {
		task := NewMockTask()
		task.On("exec").Return(nil)
		tasks = append(tasks, task)
	}

	return tasks
}

func generateFailedTasks(tasksCount int, duration time.Duration, errorRate uint8) []ExecutableTask {
	tasks := make([]ExecutableTask, 0, tasksCount)

	for i := 0; i < tasksCount; i++ {
		task := NewMockTaskWithDuration(duration)

		err := generateErrorWithErrorRate(errorRate)
		task.On("exec").Return(err)
		tasks = append(tasks, task)
	}

	return tasks
}

func generateFailedTasksWithRandomDuration(tasksCount int, errorRate uint8) []ExecutableTask {
	tasks := make([]ExecutableTask, 0, tasksCount)

	for i := 0; i < tasksCount; i++ {
		task := NewMockTask()

		err := generateErrorWithErrorRate(errorRate)
		task.On("exec").Return(err)
		tasks = append(tasks, task)
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

func assertMockExpectations(t *testing.T, tasks []ExecutableTask) {
	for _, task := range tasks {
		mockTask, ok := task.(*MockTask)
		if !ok {
			panic("tasks should be MockTask")
		}

		mockTask.AssertExpectations(t)
	}
}

func generateErrorWithErrorRate(errorRate uint8) error {
	if errorRate > 100 {
		errorRate = 100
	}

	var err error
	if rand.Float32() < float32(errorRate)/100.0 {
		err = errors.New("error")
	}

	return err
}
