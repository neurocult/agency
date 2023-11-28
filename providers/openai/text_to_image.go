package openai

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
)

type TextToImageParams struct {
	Model     string
	ImageSize string
}

func (p Provider) TextToImage(params TextToImageParams) *agency.Operation {
	return agency.NewOperation(func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
		reqBase64 := openai.ImageRequest{
			Prompt:         fmt.Sprintf("%s\n\n%s", cfg.Prompt, string(msg.Content)),
			Size:           params.ImageSize,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
			Model:          params.Model,
		}

		respBase64, err := p.client.CreateImage(ctx, reqBase64)
		if err != nil {
			return agency.Message{}, err
		}

		imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
		if err != nil {
			return agency.Message{}, err
		}

		return agency.Message{
			Role:    agency.AssistantRole,
			Content: imgBytes,
		}, nil
	})
}
