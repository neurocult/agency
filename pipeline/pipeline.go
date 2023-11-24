package pipeline

import (
	"context"

	"github.com/neurocult/agency/core"
)

type Pipeline struct {
	pipes []*core.Pipe
}

func New(pipes ...*core.Pipe) *Pipeline {
	return &Pipeline{
		pipes: pipes,
	}
}

// Interceptor is a function that is called after one pipe and before another.
type Interceptor func(in core.Message, out core.Message, cfg *core.PipeConfig)

func (p *Pipeline) Execute(ctx context.Context, input core.Message, interceptors ...Interceptor) (core.Message, error) {
	for _, pipe := range p.pipes {
		output, err := pipe.Execute(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, interceptor := range interceptors {
			interceptor(input, output, pipe.Config())
		}

		input = output
	}

	return input, nil
}
