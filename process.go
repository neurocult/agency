package agency

import (
	"context"
)

type Process struct {
	operations []*Operation
}

func NewProcess(operations ...*Operation) *Process {
	return &Process{
		operations: operations,
	}
}

// Interceptor is a function that is called after one pipe and before another.
type Interceptor func(in Message, out Message, cfg *OperationConfig)

func (p *Process) Execute(ctx context.Context, input Message, interceptors ...Interceptor) (Message, error) {
	for _, pipe := range p.operations {
		output, err := pipe.Execute(ctx, input)
		if err != nil {
			return Message{}, err
		}

		for _, interceptor := range interceptors {
			interceptor(input, output, pipe.Config())
		}

		input = output
	}

	return input, nil
}
