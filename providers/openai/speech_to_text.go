package openai

import (
	"bytes"
	"context"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
)

type SpeechToTextParams struct {
	Model       string
	Temperature NullableFloat32
}

func (f Provider) SpeechToText(params SpeechToTextParams) *agency.Operation {
	return agency.NewOperation(func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
		resp, err := f.client.CreateTranscription(ctx, openai.AudioRequest{
			Model:       params.Model,
			Prompt:      cfg.Prompt,
			FilePath:    "speech.ogg", // TODO move to cfg?
			Reader:      bytes.NewReader(msg.Content),
			Temperature: getTemperature(params.Temperature),
		})
		if err != nil {
			return agency.Message{}, err
		}

		return agency.Message{
			Role:    agency.AssistantRole,
			Content: []byte(resp.Text),
		}, nil
	})
}
