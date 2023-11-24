package openai

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
)

type ImageToTextParams struct {
	MaxTokens int
}

func (f *Factory) ImageToText(ctx context.Context, params ImageToTextParams) *agency.Operation {
	return agency.NewOperation(func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
		openaiMsg := openai.ChatCompletionMessage{
			Role:         openai.ChatMessageRoleUser,
			MultiContent: make([]openai.ChatMessagePart, 0, len(cfg.Messages)+1),
		}
		openaiMsg.MultiContent = append(openaiMsg.MultiContent, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeText,
			Text: msg.String(),
		})
		for _, msg := range cfg.Messages {
			imgBase64Str := base64.StdEncoding.EncodeToString(msg.Content)
			openaiMsg.MultiContent = append(openaiMsg.MultiContent, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL:    fmt.Sprintf("data:image/jpeg;base64,%s", imgBase64Str),
					Detail: "", // TODO
				},
			})
		}

		resp, err := f.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			MaxTokens: params.MaxTokens,
			Model:     openai.GPT4VisionPreview,
			Messages:  []openai.ChatCompletionMessage{openaiMsg},
		})
		if err != nil {
			return agency.Message{}, nil
		}

		if len(resp.Choices) < 1 {
			return agency.Message{}, errors.New("no choice")
		}
		choice := resp.Choices[0].Message // TODO what about other choices?

		return agency.Message{
			Role:    agency.AssistantRole,
			Content: []byte(choice.Content),
		}, nil
	})
}
