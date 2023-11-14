package openai

import (
	"context"
	"encoding/base64"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

func TextToImage(client *openai.Client, model string) core.PipeFactory[core.ImageConfig] {
	factory := func(options ...core.Option[core.ImageConfig]) core.Pipe {
		cfg := core.NewImageConfig(options...)

		pipe := func(ctx context.Context, msg core.Message) (core.Message, error) {
			reqBase64 := openai.ImageRequest{
				Prompt:         string(msg.Bytes()),
				Size:           cfg.Specific.Size,
				ResponseFormat: openai.CreateImageResponseFormatB64JSON,
				N:              1,
				Model:          model,
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

		return pipe
	}

	return factory
}
