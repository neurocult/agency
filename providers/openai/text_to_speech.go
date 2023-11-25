package openai

import (
	"context"
	"io"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
)

type TextToSpeechParams struct {
	Model          string
	ResponseFormat string
	Speed          float64
	Voice          string
}

func (f Provider) TextToSpeech(params TextToSpeechParams) *agency.Operation {
	return agency.NewOperation(func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
		resp, err := f.client.CreateSpeech(ctx, openai.CreateSpeechRequest{
			Model:          openai.SpeechModel(params.Model),
			Input:          msg.String(),
			Voice:          openai.SpeechVoice(params.Voice),
			ResponseFormat: openai.SpeechResponseFormat(params.ResponseFormat),
			Speed:          params.Speed,
		})
		if err != nil {
			return agency.Message{}, err
		}

		bb, err := io.ReadAll(resp)
		if err != nil {
			return agency.Message{}, err
		}

		return agency.Message{
			Role:    agency.AssistantRole,
			Content: bb,
		}, nil
	})
}
