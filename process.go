package agency

import (
	"context"
	"errors"
	"fmt"
)

// Process is a sequential chain of steps operations that can be executed in sequence.
type Process struct {
	steps []ProcessStep
}

// ProcessStep is an object that can be chained with other steps forming the process.
type ProcessStep struct {
	// Operation that current step depends on.
	// It's execution is deferred until the process reaches the corresponding step.
	Operation *Operation
	// ConfigFunc allows to modify config based a on results from the previous steps.
	// It's execution is deferred until the process reaches the corresponding step.
	ConfigFunc func(ProcessHistory, *OperationConfig) error
}

// NewProcess creates new process based on a given steps. If you don't need history use ProcessFromOperations instead.
func NewProcess(steps ...ProcessStep) *Process {
	return &Process{steps: steps}
}

// ProcessFromOperations allows to create process from operations.
// It's handy when all you need is to chain some operations together and you don't want to have an access to history.
func ProcessFromOperations(operations ...*Operation) *Process {
	steps := make([]ProcessStep, 0, len(operations))
	for _, operation := range operations {
		steps = append(steps, ProcessStep{Operation: operation, ConfigFunc: nil})
	}
	return &Process{steps: steps}
}

// ProcessInterceptor is a function that is called by Process after one step finished but before next one is started.
// Note that there's no way to modify these arguments because they relates to an operation that is already executed.
type ProcessInterceptor func(in Message, out Message, cfg OperationConfig, stepIndex uint)

// ProcessHistory stores results of the previous steps of the process. It's a process's execution context.
type ProcessHistory interface {
	Get(stepIndex uint) (Message, error) // Get takes index (starts from zero) of the step which result we want to get
	All() []Message                      // All allows to retrieve all the history of the previously processed steps
}

// processHistory implements ProcessHistory interfaces via simple slice of messages
type processHistory []Message

// Get is a panic-free way to get a message by index of the step. Indexes starts with zero. Index must be < steps count
func (p processHistory) Get(stepIndex uint) (Message, error) {
	i := int(stepIndex)
	if i >= len(p) {
		return Message{}, errors.New("step index must less than the number of steps")
	}
	return p[i], nil
}

// All simply returns p as it is.
func (p processHistory) All() []Message {
	return p
}

// Execute loops over process steps and sequentially executes them by passing output of one step as an input to another.
// If interceptors are provided, they are called on each step. So for N steps and M interceptors there's N x M executions.
func (p *Process) Execute(ctx context.Context, input Message, interceptors ...ProcessInterceptor) (Message, ProcessHistory, error) {
	history := make(processHistory, 0, len(p.steps))

	for i, step := range p.steps {
		if step.ConfigFunc != nil {
			if err := step.ConfigFunc(history, step.Operation.config); err != nil {
				return Message{}, nil, fmt.Errorf("config func on step %d: %w", i, err)
			}
		}

		output, err := step.Operation.Execute(ctx, input)
		if err != nil {
			return Message{}, nil, fmt.Errorf("operation execute: %w", err)
		}

		history = append(history, output)

		for _, interceptor := range interceptors {
			interceptor(input, output, *step.Operation.config, uint(i))
		}

		input = output
	}

	return input, history, nil
}
