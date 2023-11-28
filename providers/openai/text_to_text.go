package openai

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
)

type TextToTextParams struct {
	Model       string
	Temperature NullableFloat32
	MaxTokens   int
}

func (p Provider) TextToText(params TextToTextParams) *agency.Operation {
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

		resp, err := p.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:       params.Model,
				Temperature: getTemperature(params.Temperature),
				MaxTokens:   params.MaxTokens,
				Messages:    openAIMessages,
			},
		)
		if err != nil {
			return agency.Message{}, err
		}

		if len(resp.Choices) < 1 {
			return agency.Message{}, errors.New("no choice")
		}
		choice := resp.Choices[0].Message // TODO what about other choices?

		return agency.Message{
			Role:    agency.Role(choice.Role),
			Content: []byte(choice.Content),
		}, nil
	})
}
