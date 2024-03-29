package openai

import (
	"context"
	"errors"
	"io"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
)

type TextToStreamParams struct {
	Model       string
	Temperature NullableFloat32
	MaxTokens   int
}

type streamHandler func(delta string) error

func (p Provider) TextToStream(params TextToStreamParams, handler streamHandler) *agency.Operation {
	return agency.NewOperation(func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
		openAIMessages := make([]openai.ChatCompletionMessage, 0, len(cfg.Messages)+2)

		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: cfg.Prompt,
		})

		for _, textMsg := range cfg.Messages {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    string(textMsg.Role),
				Content: string(textMsg.Content),
			})
		}

		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: msg.String(),
		})

		resp, err := p.client.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:       params.Model,
				Temperature: getTemperature(params.Temperature),
				MaxTokens:   params.MaxTokens,
				Messages:    openAIMessages,
				Stream:      true,
			},
		)
		if err != nil {
			return agency.Message{}, err
		}
		defer resp.Close()

		var content string
		for {
			response, err := resp.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				return agency.Message{}, nil
			}

			if err := handler(response.Choices[0].Delta.Content); err != nil {
				return agency.Message{}, err
			}

			content += response.Choices[0].Delta.Content
		}

		return agency.Message{
			Role:    agency.AssistantRole,
			Content: []byte(content),
		}, nil
	})
}
