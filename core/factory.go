package core

type TextToTextFactory[T any] func(client T, params TextToTextParams) Pipe

type TextToTextParams struct {
	Model       string
	Messages    []TextMessage
	Temperature float32
}

type SpeechToTextParams struct {
	Model  string
	Prompt string
}

type SpeechToTextFactory[T any] func(client T, params SpeechToTextParams) Pipe

type TextToImageFactory[T any] func(client T, params TextToImageParams) Pipe

type TextToImageParams struct {
	Model     string
	ImageSize string
}
