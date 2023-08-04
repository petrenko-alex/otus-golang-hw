package hw06pipelineexecution

import "time"

const DataWaitLimit = time.Second * 10

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	close(out) // tmp to fail tests

	// call stages

	return out
}
