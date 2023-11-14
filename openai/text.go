package openai

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
)

var TextToText core.TextToTextFactory[*openai.Client] = func(client *openai.Client, params core.TextToTextParams) core.Pipe {
	openAIMessages := textMessagesToOpenAI(params.Messages)

	return func(ctx context.Context, msg core.Message) (core.Message, error) {
		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: string(msg.Bytes()),
		})

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

func textMessagesToOpenAI(msgs []core.TextMessage) []openai.ChatCompletionMessage {
	openAIMessages := make([]openai.ChatCompletionMessage, 0, len(msgs))
	for _, msg := range msgs {
		openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}
	return openAIMessages
}
