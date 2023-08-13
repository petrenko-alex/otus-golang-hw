package hw06pipelineexecution

import (
	"time"
)

const DataWaitLimit = time.Second * 3

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageChan := mergeDoneAndIn(done, in) // 1st stage IN chan = pipeline IN chan

	for _, stage := range stages {
		stageChan = mergeDoneAndIn(done, stage(stageChan)) // stage OUT chan = next stage IN chan
	}

	return stageChan // last stage OUT chan = pipeline OUT chan
}

func mergeDoneAndIn(done In, in In) In {
	out := make(Bi)

	go func() {
		defer close(out)
		start := time.Now()

		for {
			select {
			case val, ok := <-in:
				if !ok {
					return
				}
				start = time.Now() // reset wait timer
				out <- val
			case <-done:
				return
			default:
				spent := time.Since(start)
				if spent > DataWaitLimit {
					return
				}
			}
		}
	}()

	return out
}
