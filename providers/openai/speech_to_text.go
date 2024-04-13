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

// SpeechToText is an operation builder that creates operation than can convert speech to text.
func (f Provider) SpeechToText(params SpeechToTextParams) *agency.Operation {
	return agency.NewOperation(
		func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
			resp, err := f.client.CreateTranscription(ctx, openai.AudioRequest{
				Model:       params.Model,
				Prompt:      cfg.Prompt,
				FilePath:    "speech.ogg",
				Reader:      bytes.NewReader(msg.Content),
				Temperature: nullableToFloat32(params.Temperature),
			})
			if err != nil {
				return agency.Message{}, err
			}

			return agency.Message{
				Role:    agency.AssistantRole,
				Content: []byte(resp.Text),
			}, nil
		},
	)
}
