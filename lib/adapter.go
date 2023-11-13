package lib

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/sashabaranov/go-openai"
)

func TextPipe(client *openai.Client, prefix ...TextMessage) Pipe {
	openAIMessages := make([]openai.ChatCompletionMessage, 0, len(prefix))
	for _, textMsg := range prefix {
		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    string(textMsg.Role),
			Content: textMsg.Content,
		})
	}

	return func(ctx context.Context, msg Message) (Message, error) {
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

		return TextMessage{
			Content: choice.Content,
			Role:    Role(choice.Role),
		}, nil
	}
}

func ImagePipe(client *openai.Client) Pipe {
	return func(ctx context.Context, msg Message) (Message, error) {
		// Sample image by link
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


		return ImageMessage{
			bb: imgBytes,
		}, nil
	}
}
