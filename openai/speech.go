package openai

import (
	"bytes"
	"context"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

var SpeechToText core.SpeechToTextFactory[*openai.Client] = func(client *openai.Client, params core.SpeechToTextParams) core.Pipe {
	return func(ctx context.Context, msg core.Message) (core.Message, error) {
		resp, err := client.CreateTranscription(ctx, openai.AudioRequest{
			Model:    params.Model,
			Prompt:   params.Prompt,
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
