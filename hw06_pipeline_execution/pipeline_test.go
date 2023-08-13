package hw06pipelineexecution_test

import (
	"strconv"
	"testing"
	"time"

	. "github.com/petrenko-alex/otus-golang-hw/hw06_pipeline_execution"
	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)

				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("no stages", func(t *testing.T) {
		result := make([]int, 0, 10)
		data, inChannel := generateDataAndSendToChannel(5, nil)

		for s := range ExecutePipeline(inChannel, nil) {
			result = append(result, s.(int))
		}

		require.Equal(t, data, result)
	})

	t.Run("no data", func(t *testing.T) {
		_, inChannel := generateDataAndSendToChannel(0, nil)

		start := time.Now()
		_, opened := <-ExecutePipeline(inChannel, nil, stages...)
		elapsed := time.Since(start)

		require.False(t, opened)
		require.GreaterOrEqual(t, elapsed, DataWaitLimit)
	})

	t.Run("no data, no stages", func(t *testing.T) {
		_, inChannel := generateDataAndSendToChannel(0, nil)

		start := time.Now()
		_, opened := <-ExecutePipeline(inChannel, nil)
		elapsed := time.Since(start)

		require.False(t, opened)
		require.GreaterOrEqual(t, elapsed, DataWaitLimit)
	})

	t.Run("one stage", func(t *testing.T) {
		result := make([]int, 0, 10)
		_, inChannel := generateDataAndSendToChannel(5, nil)

		for s := range ExecutePipeline(inChannel, nil, stages[1]) {
			result = append(result, s.(int))
		}

		require.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("one element data", func(t *testing.T) {
		result := make([]string, 0, 10)
		_, inChannel := generateDataAndSendToChannel(1, nil)

		for s := range ExecutePipeline(inChannel, nil, stages...) {
			result = append(result, s.(string))
		}

		require.Equal(t, []string{"102"}, result)
	})

	t.Run("many stages", func(t *testing.T) {
		result := make([]string, 0, 10)
		data, inChannel := generateDataAndSendToChannel(5, nil)

		start := time.Now()
		for s := range ExecutePipeline(inChannel, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault),
		)
	})

	t.Run("done with no result", func(t *testing.T) {
		done := make(Bi)
		result := make([]string, 0, 10)
		_, inChannel := generateDataAndSendToChannel(5, done)

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		start := time.Now()
		for s := range ExecutePipeline(inChannel, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})

	t.Run("done with part result", func(t *testing.T) {
		done := make(Bi)
		result := make([]string, 0, 10)
		data, inChannel := generateDataAndSendToChannel(5, done)

		// Abort after all stages completed for at least 1 value + some gap
		abortDur := int(sleepPerStage)*len(stages) + int(sleepPerStage)
		go func() {
			<-time.After(time.Duration(abortDur))
			close(done)
		}()

		for s := range ExecutePipeline(inChannel, done, stages...) {
			result = append(result, s.(string))
		}

		require.NotZero(t, len(result))
		require.Less(t, len(result), len(data))
	})
}

// Генерирует слайс с тестовыми данными и отправляет их в канал.
// Возвращает тестовые данные и канал, в который производится отправка.
func generateDataAndSendToChannel(dataSize int, done In) ([]int, In) {
	in := make(Bi)
	data := make([]int, 0, dataSize)
	for i := 1; i <= dataSize; i++ {
		data = append(data, i)
	}

	if len(data) > 0 {
		go func() {
			defer close(in)

			for _, v := range data {
				select {
				case in <- v:
				case <-done:
					return
				}
			}
		}()
	}

	return data, in
}
