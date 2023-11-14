package openai

import (
	"bytes"
	"context"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

func SpeechToText(client *openai.Client, model string) core.PipeFactory[core.SpeechConfig] {
	return func(options ...core.Option[core.SpeechConfig]) core.Pipe {
		cfg := &core.Config[core.SpeechConfig]{}
		cfg.Apply(options...)

		pipe := func(ctx context.Context, msg core.Message) (core.Message, error) {
			resp, err := client.CreateTranscription(ctx, openai.AudioRequest{
				Model:    openai.Whisper1,
				FilePath: "voice.ogg",
				Reader:   bytes.NewReader(msg.Bytes()),
				Prompt:   cfg.Specific.Prompt,
			})

			if err != nil {
				return nil, err
			}

			return core.TextMessage{
				Role:    core.AssistantRole,
				Content: resp.Text,
			}, nil
		}

		return pipe
	}
}
