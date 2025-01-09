package openai

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
)

// TextToTextParams represents parameters that are specific for this operation.
type TextToTextParams struct {
	Model               string
	Temperature         NullableFloat32
	MaxTokens           int
	FuncDefs            []FuncDef
	Seed                *int
	IsToolsCallRequired bool
	Format              *openai.ChatCompletionResponseFormat
}

func (p TextToTextParams) ToolCallRequired() *string {
	var toolChoice *string
	if p.IsToolsCallRequired {
		v := "required"
		toolChoice = &v
	}

	return toolChoice
}

// TextToText is an operation builder that creates operation than can convert text to text.
// It can also call provided functions if needed, as many times as needed until the final answer is generated.
func (p Provider) TextToText(params TextToTextParams) *agency.Operation {
	openAITools := castFuncDefsToOpenAITools(params.FuncDefs)

	return agency.NewOperation(
		func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
			openAIMessages, err := agencyToOpenAIMessages(cfg, msg)
			if err != nil {
				return nil, fmt.Errorf("text to stream: %w", err)
			}

			for {
				openAIResponse, err := p.client.CreateChatCompletion(
					ctx,
					openai.ChatCompletionRequest{
						Model:          params.Model,
						Temperature:    nullableToFloat32(params.Temperature),
						MaxTokens:      params.MaxTokens,
						Messages:       openAIMessages,
						Tools:          openAITools,
						Seed:           params.Seed,
						ToolChoice:     params.ToolCallRequired(),
						ResponseFormat: params.Format,
					},
				)
				if err != nil {
					return nil, err
				}

				if len(openAIResponse.Choices) == 0 {
					return nil, errors.New("get text to text response: no choice")
				}

				responseMessage := openAIResponse.Choices[0].Message

				if len(responseMessage.ToolCalls) == 0 {
					return agency.NewTextMessage(
						agency.Role(responseMessage.Role),
						responseMessage.Content,
					), nil
				}

				openAIMessages = append(openAIMessages, responseMessage)
				for _, call := range responseMessage.ToolCalls {
					toolResponse, err := callTool(ctx, call, params.FuncDefs)
					if err != nil {
						return nil, fmt.Errorf("text to text call tool: %w", err)
					}

					if toolResponse.Role() != agency.ToolRole {
						return toolResponse, nil
					}

					openAIMessages = append(openAIMessages, toolMessageToOpenAI(toolResponse, call.ID))
				}
			}
		},
	)
}
