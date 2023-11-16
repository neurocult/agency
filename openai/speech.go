package openai

import (
	"bytes"
	"context"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

type SpeechToTextParams struct {
	Model string
}

func SpeechToText(client *openai.Client, params SpeechToTextParams) core.Pipe {
	return func(ctx context.Context, msg core.Message, options ...core.PipeOption) (core.Message, error) {
		cfg := core.NewPipeConfig(options...)

		resp, err := client.CreateTranscription(ctx, openai.AudioRequest{
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
	}
}
