package openai

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

func TextToText(client *openai.Client, model string) core.PipeFactory[core.TextConfig] {
	return func(options ...core.Option[core.TextConfig]) core.Pipe {
		cfg := core.NewTextConfig(options...)

		openAIMessages := make([]openai.ChatCompletionMessage, 0, len(cfg.Specific.Messages))
		for _, msg := range cfg.Specific.Messages {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    string(msg.Role),
				Content: msg.Content,
			})
		}

		pipe := func(ctx context.Context, msg core.Message) (core.Message, error) {
			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: string(msg.Bytes()),
			})

			resp, err := client.CreateChatCompletion(
				ctx,
				openai.ChatCompletionRequest{
					Model:       openai.GPT3Dot5Turbo,
					Messages:    openAIMessages,
					Temperature: cfg.Temperature,
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

		return pipe
	}
}
