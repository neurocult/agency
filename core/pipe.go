package core

import (
	"context"
)

// Pipe is basic building block. Pipes can be composed together into pipeline via `Then` method
type Pipe struct {
	handler      func(context.Context, Message, ...PipeOption) (Message, error)
	interceptors []Interceptor
	options      []PipeOption
}

func NewPipe(handler func(context.Context, Message, ...PipeOption) (Message, error)) *Pipe {
	return &Pipe{
		handler: handler,
	}
}

// Intercept allows execute code on each step of the pipeline.
// Interceptor called inside `Then` so it only works for pipelines with >= 2 steps
func (p *Pipe) Intercept(interceptor ...Interceptor) *Pipe {
	p.interceptors = append(p.interceptors, interceptor...)
	return p
}

// Execute executes the whole pipeline. It's just sugar over regular function call
func (p *Pipe) Execute(ctx context.Context, input Message) (Message, error) {
	output, err := p.handler(ctx, input, p.options...)
	if err != nil {
		return nil, err
	}

	for _, interceptor := range p.interceptors {
		interceptor(input, output)
	}

	return output, nil
}

// WithOptions allows to specify pipe options without execution.
func (p *Pipe) WithOptions(options ...PipeOption) *Pipe {
	p.options = options
	return p
}
