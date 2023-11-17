package core

import (
	"context"
)

// Pipe is basic building block. Pipes can be composed together into pipeline via `Then` method
type Pipe struct {
	handler      func(context.Context, Message, ...PipeOption) (Message, error)
	interceptors []func(in Message, out Message, opts ...PipeOption)
}

func NewPipe(handler func(context.Context, Message, ...PipeOption) (Message, error)) *Pipe {
	return &Pipe{
		handler: handler,
	}
}

// Intercept allows execute code on each step of the pipeline.
// Interceptor called inside `Then` so it only works for pipelines with >= 2 steps
func (p *Pipe) Intercept(interceptor ...func(Message, Message, ...PipeOption)) *Pipe {
	p.interceptors = interceptor
	return p
}

// Then takes a `next` pipe and returns new pipe that wraps `next`
func (p *Pipe) Then(next Pipe) *Pipe {
	return &Pipe{
		interceptors: p.interceptors,
		handler: func(ctx context.Context, input Message, options ...PipeOption) (Message, error) {
			output, err := p.Execute(ctx, input)
			if err != nil {
				return nil, err
			}

			return next.Execute(ctx, output)
		},
	}
}

// Execute executes the whole pipeline. It's just sugar over regular function call
func (p *Pipe) Execute(ctx context.Context, input Message) (Message, error) {
	output, err := p.handler(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, interceptor := range p.interceptors {
		interceptor(input, output)
	}

	return output, nil
}

// WithOptions allows to specify pipe options without execution.
func (p *Pipe) WithOptions(options ...PipeOption) Pipe {
	return Pipe{
		handler: func(ctx context.Context, msg Message, _ ...PipeOption) (Message, error) {
			return p.handler(ctx, msg, options...)
		},
		interceptors: p.interceptors,
	}
}
