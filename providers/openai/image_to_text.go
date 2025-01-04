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
	Model            string
	MaxTokens        int
	Temperature      NullableFloat32
	TopP             NullableFloat32
	FrequencyPenalty NullableFloat32
	PresencePenalty  NullableFloat32
}

// ImageToText is an operation builder that creates operation than can convert image to text.
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

		for _, msgFromCfg := range cfg.Messages {
			openaiMsg.MultiContent = append(
				openaiMsg.MultiContent,
				bytesToOpenAIMessagePart(msgFromCfg),
			)
		}

		openaiMsg.MultiContent = append(
			openaiMsg.MultiContent,
			bytesToOpenAIMessagePart(msg),
		)

		resp, err := f.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			MaxTokens:        params.MaxTokens,
			Model:            params.Model,
			Messages:         []openai.ChatCompletionMessage{openaiMsg},
			Temperature:      nullableToFloat32(params.Temperature),
			TopP:             nullableToFloat32(params.TopP),
			FrequencyPenalty: nullableToFloat32(params.FrequencyPenalty),
			PresencePenalty:  nullableToFloat32(params.PresencePenalty),
		})
		if err != nil {
			return nil, err
		}

		if len(resp.Choices) < 1 {
			return nil, errors.New("no choice")
		}
		choice := resp.Choices[0].Message

		return agency.NewTextMessage(agency.AssistantRole, choice.Content), nil
	})
}

func bytesToOpenAIMessagePart(msg agency.Message) openai.ChatMessagePart {
	imgMsg := msg.(agency.ImageMessage) // panic if given msg is not image

	return openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeImageURL,
		ImageURL: &openai.ChatMessageImageURL{
			URL: fmt.Sprintf(
				"data:image/jpeg;base64,%s",
				base64.StdEncoding.EncodeToString(imgMsg.Content()),
			),
			Detail: openai.ImageURLDetailAuto,
		},
		Text: imgMsg.Description(),
	}
}
