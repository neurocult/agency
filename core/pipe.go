package core

import (
	"context"
)

type PipeHandler func(context.Context, Message, *PipeConfig) (Message, error)

// Pipe is basic building block.
type Pipe struct {
	handler PipeHandler
	config  *PipeConfig
}

func (p *Pipe) Config() *PipeConfig {
	return p.config
}

// PipeConfig represents abstract pipe configuration.
// It contains fields for all possible modalities but nothing specific to concrete model implementations.
type PipeConfig struct {
	Prompt   string
	Messages []Message
}

// NewPipe allows to create Pipe from a function.
func NewPipe(handler PipeHandler) *Pipe {
	return &Pipe{
		handler: handler,
		config:  &PipeConfig{},
	}
}

// Execute executes the whole pipeline.
func (p *Pipe) Execute(ctx context.Context, input Message) (Message, error) {
	output, err := p.handler(ctx, input, p.config)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (p *Pipe) SetPrompt(prompt string, args ...any) *Pipe {
	p.config.Prompt = prompt
	return p
}

func (p *Pipe) SetMessages(msgs []Message) *Pipe {
	p.config.Messages = msgs
	return p
}
