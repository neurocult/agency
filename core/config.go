package core

type Config[T SpecificConfig] struct {
	Temperature float32
	Specific    T
}

func (c *Config[SpecificConfig]) Apply(options ...Option[SpecificConfig]) {
	for _, option := range options {
		option(c)
	}
}

type SpecificConfig interface {
	TextConfig | ImageConfig | SpeechConfig
}

type TextConfig struct {
	Messages []TextMessage
}

func NewTextConfig(options ...Option[TextConfig]) *Config[TextConfig] {
	cfg := &Config[TextConfig]{}
	cfg.Apply(options...)
	return cfg
}

type ImageConfig struct {
	Size string
}

func WithSize(size string) Option[ImageConfig] {
	return func(cfg *Config[ImageConfig]) {
		cfg.Specific.Size = size
	}
}

func NewImageConfig(options ...Option[ImageConfig]) *Config[ImageConfig] {
	cfg := &Config[ImageConfig]{}
	cfg.Apply(options...)
	return cfg
}

type SpeechConfig struct {
	Prompt string
}

func NewSpeechConfig(options ...Option[SpeechConfig]) *Config[SpeechConfig] {
	cfg := &Config[SpeechConfig]{}
	cfg.Apply(options...)
	return cfg
}

func WithPrompt(prompt string) Option[SpeechConfig] {
	return func(cfg *Config[SpeechConfig]) {
		cfg.Specific.Prompt = prompt
	}
}

type PipeFactory[T SpecificConfig] func(options ...Option[T]) Pipe

// Pipe is just a sintactic sugar over regular function call
func (p PipeFactory[SpecificConfig]) Pipe(options ...Option[SpecificConfig]) Pipe {
	return p(options...)
}

type Option[T SpecificConfig] func(cfg *Config[T])

func WithTemperature[T SpecificConfig](temperature float32) Option[T] {
	return func(cfg *Config[T]) {
		cfg.Temperature = temperature
	}
}

func WithSystemMessage(prompt string, args ...any) Option[TextConfig] {
	return func(cfg *Config[TextConfig]) {
		systemMsg := NewSystemMessage(prompt).Bind(args...)
		cfg.Specific.Messages = append(cfg.Specific.Messages, systemMsg)
	}
}
