package openai

import (
	"bytes"
	"context"

	"github.com/eqtlab/lib/core"
	"github.com/sashabaranov/go-openai"
)

type SpeechToTextParams struct {
	Model string
}

func (f Factory) SpeechToText(params SpeechToTextParams) *core.Pipe {
	return core.NewPipe(func(ctx context.Context, msg core.Message, options ...core.PipeOption) (core.Message, error) {
		cfg := core.NewPipeConfig(options...)

		resp, err := f.client.CreateTranscription(ctx, openai.AudioRequest{
			Model:    params.Model,
			Prompt:   cfg.Prompt,
			FilePath: "voice.ogg",
			Reader:   bytes.NewReader(msg.Bytes()),
		})
		if err != nil {
			return nil, err
		}

		return core.TextMessage{
			Role:    core.AssistantRole,
			Content: resp.Text,
		}, nil
	})
}
