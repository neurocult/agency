package agency

import (
	"context"
)

// Process is a chain of operations that can be executed in sequence.
type Process struct {
	operations []*Operation
}

func NewProcess(operations ...*Operation) *Process {
	return &Process{
		operations: operations,
	}
}

// Interceptor is a function that is called by Process after one operation finished but before next one is started.
type Interceptor func(in Message, out Message, cfg *OperationConfig)

// Execute iterates over Process's operations and sequentially executes them.
// After first operation is executed it uses its output as an input to the second one and so on until the whole chain is finished.
// It also executes all given interceptors, if they are provided, so for every N operations and M interceptors it's N x M executions.
func (p *Process) Execute(ctx context.Context, input Message, interceptors ...Interceptor) (Message, error) {
	for _, operation := range p.operations {
		output, err := operation.Execute(ctx, input)
		if err != nil {
			return Message{}, err
		}

		// FIXME while these are called AFTER operation and not before it's impossible to modify configuration
		for _, interceptor := range interceptors {
			interceptor(input, output, operation.Config())
		}

		input = output
	}

	return input, nil
}
