package openai

import (
	"context"
	"errors"
	"fmt"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// FuncDef represents a function definition that can be called during the conversation.
type FuncDef struct {
	Name        string
	Description string
	// Parameters is an optional structure that defines the schema of the parameters that the function accepts.
	Parameters *jsonschema.Definition
	// Body is the actual function that get's called.
	// Parameters passed are bytes that can be unmarshalled to type that implements provided json schema.
	// Returned result must be anything that can be marshalled, including primitive values.
	Body func(ctx context.Context, params []byte) (agency.Message, error)
}

type FuncDefs []FuncDef

func (ds FuncDefs) getFuncDefByName(name string) *FuncDef {
	for _, f := range ds {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

func toolMessageToOpenAI(message agency.Message, toolID string) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:       string(message.Role()),
		Content:    string(message.Content()),
		ToolCallID: toolID,
	}
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
		}
		tools = append(tools, tool)
	}
	return tools
}
