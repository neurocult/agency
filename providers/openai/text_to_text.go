package openai

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"

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
			openAIMessages, err := agencyToOpenaiMessages(cfg, msg)
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
					return OpenaiToAgencyMessage(responseMessage), nil
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

// === Helpers ===

func castFuncDefsToOpenAITools(funcDefs []FuncDef) []openai.Tool {
	tools := make([]openai.Tool, 0, len(funcDefs))
	for _, f := range funcDefs {
		tool := openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
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

func agencyToOpenaiMessages(cfg *agency.OperationConfig, msg agency.Message) ([]openai.ChatCompletionMessage, error) {
	openAIMessages := make([]openai.ChatCompletionMessage, 0, len(cfg.Messages)+2)

	openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: cfg.Prompt,
	})

	for _, cfgMsg := range cfg.Messages {
		openAIMessages = append(openAIMessages, messageToOpenAI(cfgMsg))
	}

	openaiMsg := openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
	}

	switch msg.Kind() {
	case agency.TextKind:
		openaiMsg.Content = string(msg.Content())
	case agency.ImageKind:
		openaiMsg.MultiContent = append(
			openaiMsg.MultiContent,
			openAIBase64ImageMessage(msg.Content()),
		)
	default:
		return nil, fmt.Errorf("operator doesn't support %s kind", msg.Kind())
	}

	openAIMessages = append(openAIMessages, openaiMsg)

	return openAIMessages, nil
}

func callTool(
	ctx context.Context,
	call openai.ToolCall,
	defs FuncDefs,
) (agency.Message, error) {
	funcToCall := defs.getFuncDefByName(call.Function.Name)
	if funcToCall == nil {
		return nil, errors.New("function not found")
	}

	funcResult, err := funcToCall.Body(ctx, []byte(call.Function.Arguments))
	if err != nil {
		return funcResult, fmt.Errorf("call function %s: %w", funcToCall.Name, err)
	}

	return funcResult, nil
}

func OpenaiToAgencyMessage(msg openai.ChatCompletionMessage) agency.Message {
	return agency.NewTextMessage(
		agency.Role(msg.Role),
		msg.Content,
	)
}
