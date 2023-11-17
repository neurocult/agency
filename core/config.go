package core

import "fmt"

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
