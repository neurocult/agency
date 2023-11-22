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

func (p *Pipe) WithPrompt(prompt string, args ...any) *Pipe {
	p.config.Prompt = prompt
	return p
}

func (p *Pipe) WithMessages(msgs []Message) *Pipe {
	p.config.Messages = msgs
	return p
}

// pipe1 = text.WithPrompt().WithMessages()

// ---

// text1 = openai.TextToText(params...)
// text2 = openai.TextToText(params...)
// text3 = openai.TextToText(params...)

// // business logic
// translate = text1.WithPrompt(1).WithMessages([""])
// uppercase = text2.WithPrompt(2).WithMessages([""])
// replaceWhitespaces = text3.WithPrompt(3).WithMessages([""])
