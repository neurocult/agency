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
	Temperature NullableFloat32
	MaxTokens   int
}

func (f *Provider) ImageToText(params ImageToTextParams) *agency.Operation {
	return agency.NewOperation(func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
		openaiMsg := openai.ChatCompletionMessage{
			Role:         openai.ChatMessageRoleUser,
			MultiContent: make([]openai.ChatMessagePart, 0, len(cfg.Messages)+2),
		}

		openaiMsg.MultiContent = append(openaiMsg.MultiContent, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeText,
			Text: cfg.Prompt,
		})

		for _, cfgMsg := range cfg.Messages {
			openaiMsg.MultiContent = append(
				openaiMsg.MultiContent,
				openAIBase64ImageMessage(cfgMsg.Content),
			)
		}

		openaiMsg.MultiContent = append(
			openaiMsg.MultiContent,
			openAIBase64ImageMessage(msg.Content),
		)

		resp, err := f.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			MaxTokens:   params.MaxTokens,
			Model:       openai.GPT4VisionPreview,
			Messages:    []openai.ChatCompletionMessage{openaiMsg},
			Temperature: getTemperature(params.Temperature),
		})
		if err != nil {
			return agency.Message{}, err
		}

		if len(resp.Choices) < 1 {
			return agency.Message{}, errors.New("no choice")
		}
		choice := resp.Choices[0].Message

		return agency.Message{
			Role:    agency.AssistantRole,
			Content: []byte(choice.Content),
		}, nil
	})
}

func openAIBase64ImageMessage(bb []byte) openai.ChatMessagePart {
	imgBase64Str := base64.StdEncoding.EncodeToString(bb)
	return openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeImageURL,
		ImageURL: &openai.ChatMessageImageURL{
			URL:    fmt.Sprintf("data:image/jpeg;base64,%s", imgBase64Str),
			Detail: openai.ImageURLDetailAuto,
		},
	}
}
