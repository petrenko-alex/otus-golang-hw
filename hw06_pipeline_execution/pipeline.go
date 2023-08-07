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
	stageChan := in // 1st stage IN chan = pipeline IN chan

	for _, stage := range stages {
		stageChan = stage(stageChan) // stage OUT chan = next stage IN chan
	}

	return stageChan // last stage OUT chan = pipeline OUT chan
}
