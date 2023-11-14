package openai

import (
	"context"
	"errors"

	"github.com/eqtlab/lib/core"
	"github.com/sashabaranov/go-openai"
)

func NewTextToText(client *openai.Client, model string) core.Configurator[core.TextConfig] {
	return func(options ...core.ConfiguratorOption[core.TextConfig]) core.Pipe {
		cfg := &core.Config[core.TextConfig]{}

		for _, opt := range options {
			opt(cfg)
		}

		openAIMessages := make([]openai.ChatCompletionMessage, 0, len(options))
		for _, msg := range cfg.Model.Messages {
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
