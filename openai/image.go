package openai

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

type TextToImageParams struct {
	Model     string
	ImageSize string
}

func TextToImage(client *openai.Client, params TextToImageParams) core.Pipe {
	return func(ctx context.Context, msg core.Message, options ...core.PipeOption) (core.Message, error) {
		cfg := core.NewPipeConfig(options...)

		reqBase64 := openai.ImageRequest{
			Prompt:         fmt.Sprintf("%s\n\n%s", cfg.Prompt, string(msg.Bytes())),
			Size:           params.ImageSize,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
			Model:          params.Model,
		}

		respBase64, err := client.CreateImage(ctx, reqBase64)
		if err != nil {
			return nil, err
		}

		imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
		if err != nil {
			return nil, err
		}

		return core.NewImageMessage(imgBytes), nil
	}
}
