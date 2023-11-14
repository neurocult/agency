package openai

import (
	"context"
	"encoding/base64"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

var TextToImage core.TextToImageFactory[*openai.Client] = func(client *openai.Client, params core.TextToImageParams) core.Pipe {
	return func(ctx context.Context, msg core.Message) (core.Message, error) {
		reqBase64 := openai.ImageRequest{
			Prompt:         string(msg.Bytes()),
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
