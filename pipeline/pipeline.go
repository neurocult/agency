package pipeline

import (
	"context"

	"github.com/eqtlab/lib/core"
)

type Pipeline struct {
	pipes        []*core.Pipe
	interceptors []core.Interceptor
}

func New(pipes ...*core.Pipe) *Pipeline {
	return &Pipeline{
		pipes: pipes,
	}
}

func (p *Pipeline) InterceptEach(interceptor core.Interceptor) *Pipeline {
	p.interceptors = append(p.interceptors, interceptor)
	return p
}

func (p *Pipeline) Execute(ctx context.Context, message core.Message) (core.Message, error) {
	for _, pipe := range p.pipes {
		pipe.Intercept(p.interceptors...)

		var err error
		message, err = pipe.Execute(ctx, message)
		if err != nil {
			return nil, err
		}
	}

	return message, nil
}
