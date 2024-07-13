package openai

import (
	"context"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type ToolResultMessage struct {
	agency.Message

	ToolID   string
	ToolName string
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
