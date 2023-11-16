package core

import (
	"context"
)

// Pipe is basic building block. Pipes can be composed together into pipeline via `Then` method
type Pipe struct {
	handler     func(context.Context, Message, ...PipeOption) (Message, error)
	interceptor func(Message, ...PipeOption)
}

func NewPipe(handler func(context.Context, Message, ...PipeOption) (Message, error)) *Pipe {
	return &Pipe{
		handler: handler,
	}
}

// Intercept allows execute code on each step of the pipeline.
// Interceptor called inside `Then` so it only works for pipelines with >= 2 steps
func (p *Pipe) Intercept(interceptor func(Message, ...PipeOption)) *Pipe {
	p.interceptor = interceptor
	return p
}

// Then takes a `next` pipe and returns new pipe that wraps `next`
func (p *Pipe) Then(next Pipe) *Pipe {
	return &Pipe{
		interceptor: p.interceptor,
		handler: func(ctx context.Context, input Message, options ...PipeOption) (Message, error) {
			output, err := p.handler(ctx, input)
			if err != nil {
				return nil, err
			}

			// FIXME does not work for first and last call
			if p.interceptor != nil {
				p.interceptor(output, options...)
			}

			return next.handler(ctx, output)
		},
	}
}

// Execute executes the whole pipeline. It's just sugar over regular function call
func (p *Pipe) Execute(ctx context.Context, msg Message) (Message, error) {
	return p.handler(ctx, msg)
}

// WithOptions allows to specify pipe options without execution.
func (p *Pipe) WithOptions(options ...PipeOption) Pipe {
	return Pipe{
		handler: func(ctx context.Context, msg Message, _ ...PipeOption) (Message, error) {
			return p.handler(ctx, msg, options...)
		},
		interceptor: p.interceptor,
	}
}
