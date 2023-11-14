package core

import (
	"context"
)

// Pipe is basic building block. Pipes can be composed together into pipeline
type Pipe func(context.Context, Message) (Message, error)

// Then takes a `next` pipe and returns new pipe that wraps `next`
func (p Pipe) Then(next Pipe) Pipe {
	return func(ctx context.Context, bb Message) (Message, error) {
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
