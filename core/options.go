package core

type Config[T any] struct {
	Temperature float32
	Model       T
}

type TextConfig struct {
	Messages []Message
}

type Configurator[T any] func(options ...ConfiguratorOption[T]) Pipe

type ConfiguratorOption[T any] func(cfg *Config[T])

func WithTemperature[T TextConfig](temperature float32) ConfiguratorOption[T] {
	return func(cfg *Config[T]) {
		cfg.Temperature = temperature
	}
}

func WithMessages(msgs ...Message) ConfiguratorOption[TextConfig] {
	return func(cfg *Config[TextConfig]) {
		cfg.Model.Messages = msgs
	}
}

func WithPrompt(prompt string, args ...any) ConfiguratorOption[TextConfig] {
	return func(cfg *Config[TextConfig]) {
		cfg.Model.Messages = append(cfg.Model.Messages, NewSystemMessage(prompt).Bind(args...))
	}
}
