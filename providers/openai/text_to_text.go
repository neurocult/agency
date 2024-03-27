package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"

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
	// Parameters is an optional structure that defines the schema of the parameters that the function accepts.
	Parameters *jsonschema.Definition
	// Body is the actual function that get's called.
	// Parameters passed are bytes that can be unmarshalled to type that implements provided json schema.
	// Returned result must be anything that can be marshalled, including primitive values.
	Body func(ctx context.Context, params []byte) (any, error)
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
				fmt.Println("===ITERATION===")

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
				firstChoice := openAIResponse.Choices[0]

				fmt.Println(firstChoice.Message.ToolCalls)

				if len(firstChoice.Message.ToolCalls) == 0 {
					return agency.Message{
						Role:    agency.Role(firstChoice.Message.Role),
						Content: []byte(firstChoice.Message.Content),
					}, nil
				}

				firstToolCall := firstChoice.Message.ToolCalls[0]
				funcToCall := getFuncDefByName(params.FuncDefs, firstToolCall.Function.Name)
				if funcToCall == nil {
					return agency.Message{}, errors.New("function not found")
				}

				funcResult, err := funcToCall.Body(ctx, []byte(firstToolCall.Function.Arguments))
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
		tool := openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: openai.FunctionDefinition{
				Name:        f.Name,
				Description: f.Description,
			},
		}
		if f.Parameters != nil {
			tool.Function.Parameters = f.Parameters
		} else {
			tool.Function.Parameters = jsonschema.Definition{
				Type: jsonschema.Object, // because we can't pass empty parameters
			}
		}
		tools = append(tools, tool)
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
