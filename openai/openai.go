package openai

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

type PipeFactory struct {
	client *openai.Client
}

func (p PipeFactory) TextToText(prefix ...core.TextMessage) core.Pipe {
	openAIMessages := make([]openai.ChatCompletionMessage, 0, len(prefix))
	for _, textMsg := range prefix {
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

		resp, err := p.client.CreateChatCompletion(
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

func (p PipeFactory) TextToImage() core.Pipe {
	return func(ctx context.Context, msg core.Message) (core.Message, error) {
		reqBase64 := openai.ImageRequest{
			Prompt:         string(msg.Bytes()),
			Size:           openai.CreateImageSize256x256,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
		}

		respBase64, err := p.client.CreateImage(ctx, reqBase64)
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

func (p PipeFactory) SpeechToText() core.Pipe {
	return func(ctx context.Context, msg core.Message) (core.Message, error) {
		resp, err := p.client.CreateTranscription(ctx, openai.AudioRequest{
			Model:    openai.Whisper1,
			FilePath: "voice.ogg",
			Reader:   bytes.NewReader(msg.Bytes()),
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

func NewPipeFactory(client *openai.Client) PipeFactory {
	return PipeFactory{
		client: client,
	}
}
