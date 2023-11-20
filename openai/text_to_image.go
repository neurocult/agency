package openai

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/eqtlab/lib/core"
	"github.com/sashabaranov/go-openai"
)

type TextToImageParams struct {
	Model     string
	ImageSize string
}

func (f Factory) TextToImage(params TextToImageParams) *core.Pipe {
	return core.NewPipe(func(ctx context.Context, msg core.Message, options ...core.PipeOption) (core.Message, error) {
		cfg := core.NewPipeConfig(options...)

		reqBase64 := openai.ImageRequest{
			Prompt:         fmt.Sprintf("%s\n\n%s", cfg.Prompt, string(msg.Bytes())),
			Size:           params.ImageSize,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
			Model:          params.Model,
		}

		respBase64, err := f.client.CreateImage(ctx, reqBase64)
		if err != nil {
			return nil, err
		}

		imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
		if err != nil {
			return nil, err
		}

		return core.NewImageMessage(imgBytes), nil
	})
}
