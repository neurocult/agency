package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sashabaranov/go-openai/jsonschema"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	t2tOp := openai.
		New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		TextToText(openai.TextToTextParams{
			Model: "gpt-3.5-turbo",
			FuncDefs: []openai.FuncDef{
				// function without parameters
				{
					Name:        "GetMeaningOfLife",
					Description: "Answer questions about meaning of life",
					Body: func(ctx context.Context, _ []byte) (any, error) {
						return 42, nil
					},
				},
				// function with parameters
				{
					Name:        "ChangeNumbers",
					Description: "Change given numbers when asked",
					Parameters: &jsonschema.Definition{
						Type: "object",
						Properties: map[string]jsonschema.Definition{
							"a": {Type: "integer"},
							"b": {Type: "integer"},
						},
					},
					Body: func(ctx context.Context, params []byte) (any, error) {
						var pp struct{ A, B int }
						if err := json.Unmarshal(params, &pp); err != nil {
							return nil, err
						}
						return (pp.A + pp.B) * 10, nil // *10 is just to distinguish from normal response
					},
				},
			},
		}).
		SetPrompt(`
Answer questions about meaning of life and summing numbers.
Always use GetMeaningOfLife and ChangeNumbers functions results as answers.
Examples:
- User: what is the meaning of life?
- Assistant: 42
- User: 1+1
- Assistant: 20
- User: 1+1 and what is the meaning of life?
- Assistant: 20 and 42`)

	ctx := context.Background()

	// test for first function call
	answer, err := t2tOp.Execute(
		ctx,
		agency.NewMessage(agency.UserRole, agency.TextKind, []byte("what is the meaning of life?")),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)

	// test for second function call
	answer, err = t2tOp.Execute(
		ctx,
		agency.NewMessage(agency.UserRole, agency.TextKind, []byte("1+1?")),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)

	// test for both function calls at the same time
	answer, err = t2tOp.Execute(
		ctx,
		agency.NewMessage(agency.UserRole, agency.TextKind, []byte("1+1 and what is the meaning of life?")),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}
