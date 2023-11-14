package openai

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

func NewTextToText(client *openai.Client, model string) core.Configurator {
	return func(prefix ...core.Message) core.Pipe {
		openAIMessages := make([]openai.ChatCompletionMessage, 0, len(prefix))
		for _, msg := range prefix {
			textMsg, ok := msg.(core.TextMessage)
			if !ok {
				panic("not ok") // TODO handle err
			}

			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    string(textMsg.Role),
				Content: textMsg.Content,
			})
		}

		return func(ctx context.Context, msg core.Message) (core.Message, error) {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: string(msg.Bytes()),
			})

			resp, err := client.CreateChatCompletion(
				ctx,
				openai.ChatCompletionRequest{
					Model:    openai.GPT3Dot5Turbo,
					Messages: openAIMessages,
				},
			)
			if err != nil {
				return nil, err
			}

			if len(resp.Choices) < 1 {
				return nil, errors.New("no choice")
			}
			choice := resp.Choices[0].Message

			return core.TextMessage{
				Content: choice.Content,
				Role:    core.Role(choice.Role),
			}, nil
		}
	}
}

func NewTextToImage(client *openai.Client, model string) core.Configurator {
	return func(prefix ...core.Message) core.Pipe {
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

func NewSpeechToText(client *openai.Client, model string) core.Configurator {
	return func(prefix ...core.Message) core.Pipe {
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
