package core

import "fmt"

// PipeConfig represents abstract pipe configuration.
// It contains fields for all possible modalities but nothing specific to concrete model implementations.
// It allows dynamically create variations of pipes depending on request.
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
