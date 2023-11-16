package core

import (
	"context"
	"fmt"
)

type PipeConfig struct {
	Prompt string
}

type PipeOption func(*PipeConfig)

func NewPipeConfig(options ...PipeOption) *PipeConfig {
	c := &PipeConfig{}
	for _, option := range options {
		option(c)
	}

	return c
}

func WithPrompt(prompt string, args ...any) PipeOption {
	return func(c *PipeConfig) {
		c.Prompt = fmt.Sprintf(prompt, args...)
	}
}

// Pipe is basic building block. Pipes can be composed together into pipeline via `Then` method
type Pipe func(context.Context, Message, ...PipeOption) (Message, error)

// Then takes a `next` pipe and returns new pipe that wraps `next`
func (p Pipe) Then(next Pipe) Pipe {
	return func(ctx context.Context, bb Message, options ...PipeOption) (Message, error) {
		bb, err := p(ctx, bb)
		if err != nil {
			return nil, err
		}
		return next(ctx, bb)
	}
}

// Execute executes the pipe(line). This is syntactic sugar of regular function call
func (p Pipe) Execute(ctx context.Context, bb Message) (Message, error) {
	return p(ctx, bb)
}

func (p Pipe) WithOptions(options ...PipeOption) Pipe {
	return func(ctx context.Context, bb Message, _ ...PipeOption) (Message, error) {
		bb, err := p(ctx, bb)
		if err != nil {
			return nil, err
		}
		return p(ctx, bb, options...)
	}
}
