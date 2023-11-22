package openai

import (
	"context"
	"errors"

	"github.com/eqtlab/lib/core"
	"github.com/sashabaranov/go-openai"
)

type TextToTextParams struct {
	Model       string
	Temperature float32
}

func (f Factory) TextToText(params TextToTextParams) *core.Pipe {
	return core.NewPipe(func(ctx context.Context, msg core.Message, cfg *core.PipeConfig) (core.Message, error) {
		openAIMessages := make([]openai.ChatCompletionMessage, 0, len(cfg.Messages)+2)

		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: cfg.Prompt,
		})

		for _, msg = range cfg.Messages {
			textMsg, ok := msg.(core.TextMessage)
			if !ok {
				return nil, errors.New("...")
			}

			openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
				Role:    string(textMsg.Role),
				Content: textMsg.Content,
			})
		}

		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: msg.String(),
		})

		resp, err := f.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:       params.Model,
				Temperature: params.Temperature,
				Messages:    openAIMessages,
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
	})
}
