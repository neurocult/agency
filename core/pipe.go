package core

import (
	"context"
)

// Pipe is basic building block.
type Pipe struct {
	handler      func(context.Context, Message, ...PipeOption) (Message, error)
	interceptors []Interceptor
	options      []PipeOption
}

// NewPipe allows to create Pipe from a function.
func NewPipe(handler func(context.Context, Message, ...PipeOption) (Message, error)) *Pipe {
	return &Pipe{
		handler: handler,
	}
}

// After allows execute code after pipe to intercept execution between pipes.
func (p *Pipe) After(interceptor ...Interceptor) *Pipe {
	p.interceptors = append(p.interceptors, interceptor...)
	return p
}

// Execute executes the whole pipeline.
func (p *Pipe) Execute(ctx context.Context, input Message) (Message, error) {
	output, err := p.handler(ctx, input, p.options...)
	if err != nil {
		return nil, err
	}

	for _, interceptor := range p.interceptors {
		interceptor(input, output, p.options...)
	}

	return output, nil
}

// WithOptions returns new Pipe with specified options.
func (p *Pipe) WithOptions(options ...PipeOption) *Pipe {
	return &Pipe{
		handler:      p.handler,
		interceptors: p.interceptors,
		options:      options,
	}
}
