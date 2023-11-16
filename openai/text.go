package openai

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

type TextToTextParams struct {
	Model       string
	Temperature float32
}

func TextToText(client *openai.Client, params TextToTextParams) core.Pipe {
	return func(ctx context.Context, msg core.Message, options ...core.PipeOption) (core.Message, error) {
		cfg := core.NewPipeConfig(options...)

		openAIMessages := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: string(cfg.Prompt),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: string(msg.Bytes()),
			},
		}

		resp, err := client.CreateChatCompletion(
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
	}
}
