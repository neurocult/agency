package openai

import (
	"bytes"
	"context"
	"encoding/base64"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

type ImageConfig struct {
}
type ImageConfiguratorOption func(cfg *SpeechConfig)

type ImageConfigurator func(...ImageConfiguratorOption) core.Pipe

func NewTextToImage(client *openai.Client, model string) ImageConfigurator {
	return func(prefix ...ImageConfiguratorOption) core.Pipe {
		return func(ctx context.Context, msg core.Message) (core.Message, error) {
			reqBase64 := openai.ImageRequest{
				Prompt:         string(msg.Bytes()),
				Size:           openai.CreateImageSize256x256,
				ResponseFormat: openai.CreateImageResponseFormatB64JSON,
				N:              1,
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
}

type SpeechConfig struct {
}
type SpeechConfiguratorOption func(cfg *SpeechConfig)

type SpeechConfigurator func(...SpeechConfiguratorOption) core.Pipe

func NewSpeechToText(client *openai.Client, model string) SpeechConfigurator {
	return func(prefix ...SpeechConfiguratorOption) core.Pipe {
		return func(ctx context.Context, msg core.Message) (core.Message, error) {
			resp, err := client.CreateTranscription(ctx, openai.AudioRequest{
				Model:    openai.Whisper1,
				FilePath: "voice.ogg",
				Reader:   bytes.NewReader(msg.Bytes()),
				// Prompt: , TODO use prefix here
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
}
