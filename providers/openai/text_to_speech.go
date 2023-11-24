package openai

import (
	"context"
	"io"

	"github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency/core"
)

type TextToSpeechParams struct {
	Model          string
	ResponseFormat string
	Speed          float64
	Voice          string
}

func (f Factory) TextToSpeech(params TextToSpeechParams) *core.Pipe {
	return core.NewPipe(func(ctx context.Context, msg core.Message, cfg *core.PipeConfig) (core.Message, error) {
		resp, err := f.client.CreateSpeech(ctx, openai.CreateSpeechRequest{
			Model:          openai.SpeechModel(params.Model),
			Input:          msg.String(),
			Voice:          openai.SpeechVoice(params.Voice),
			ResponseFormat: openai.SpeechResponseFormat(params.ResponseFormat),
			Speed:          params.Speed,
		})
		if err != nil {
			return nil, err
		}

		bb, err := io.ReadAll(resp)
		if err != nil {
			return nil, err
		}

		return core.NewSpeechMessage(bb), nil
	})
}
