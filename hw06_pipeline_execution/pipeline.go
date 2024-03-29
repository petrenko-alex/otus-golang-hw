package hw06pipelineexecution

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

		for {
			select {
			case val, ok := <-in:
				if !ok {
					return
				}
				out <- val
			case <-done:
				return
			}
		}
	}()

	return out
}
