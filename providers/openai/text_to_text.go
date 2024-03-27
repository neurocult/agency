package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
)

// TextToTextParams represents parameters that are specific for this operation.
type TextToTextParams struct {
	Model       string
	Temperature NullableFloat32
	MaxTokens   int
	FuncDefs    []FuncDef
}

// FuncDef represents a function definition that can be called during the conversation.
type FuncDef struct {
	Name        string
	Description string
	Parameters  any // Parameters is a structure that defines the schema of the parameters that the function accepts.
	// Body is the actual function that get's called.
	// Parameters must be pointer to a structure that matches `Parameters` schema via json-tags.
	// Returned result must be json-marshallable object.
	Body func(ctx context.Context, params any) (any, error)
}

// TextToText is an operation builder that creates operation than can convert text to text.
// It can also call provided functions if needed, as many times as needed until the final answer is generated.
func (p Provider) TextToText(params TextToTextParams) *agency.Operation {
	openAITools := castFuncDefsToOpenAITools(params.FuncDefs)

	return agency.NewOperation(
		func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
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

			for {
				openAIResponse, err := p.client.CreateChatCompletion(
					ctx,
					openai.ChatCompletionRequest{
						Model:       params.Model,
						Temperature: getTemperature(params.Temperature),
						MaxTokens:   params.MaxTokens,
						Messages:    openAIMessages,
						Tools:       openAITools,
					},
				)
				if err != nil {
					return agency.Message{}, err
				}

				if len(openAIResponse.Choices) < 1 {
					return agency.Message{}, errors.New("no choice")
				}
				answer := openAIResponse.Choices[0]

				if answer.FinishReason != openai.FinishReasonFunctionCall {
					return agency.Message{
						Role:    agency.Role(answer.Message.Role),
						Content: []byte(answer.Message.Content),
					}, nil
				}

				funcToCall := getFuncDefByName(params.FuncDefs, answer.Message.FunctionCall.Name)
				if funcToCall == nil {
					return agency.Message{}, errors.New("function not found")
				}

				var params = funcToCall.Parameters
				if err = json.Unmarshal([]byte(answer.Message.FunctionCall.Arguments), &params); err != nil {
					return agency.Message{}, fmt.Errorf(
						"unmarshal %s arguments: %w",
						answer.Message.FunctionCall.Name, err,
					)
				}

				funcResult, err := funcToCall.Body(ctx, params)
				if err != nil {
					return agency.Message{}, fmt.Errorf("call function %s: %w", funcToCall.Name, err)
				}

				bb, err := json.Marshal(funcResult)
				if err != nil {
					return agency.Message{}, fmt.Errorf("marshal function result: %w", err)
				}

				openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: string(bb),
				})
			}
		},
	)
}

// === Helpers ===

func castFuncDefsToOpenAITools(funcDefs []FuncDef) []openai.Tool {
	tools := make([]openai.Tool, 0, len(funcDefs))
	for _, f := range funcDefs {
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: openai.FunctionDefinition{
				Name:        f.Name,
				Description: f.Description,
				Parameters:  f.Parameters,
			},
		})
	}
	return tools
}

func getFuncDefByName(funcDefs []FuncDef, name string) *FuncDef {
	for _, f := range funcDefs {
		if f.Name == name {
			return &f
		}
	}
	return nil
}
