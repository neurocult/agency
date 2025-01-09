package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	go_openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	t2tOp := openai.
		New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		TextToText(openai.TextToTextParams{
			Model: go_openai.GPT4oMini,
			FuncDefs: []openai.FuncDef{
				// function without parameters
				{
					Name:        "GetMeaningOfLife",
					Description: "Answer questions about meaning of life",
					Body: func(ctx context.Context, _ []byte) (agency.Message, error) {
						// because we don't need any arguments
						return agency.NewTextMessage(agency.ToolRole, "42"), nil
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
					Body: func(ctx context.Context, params []byte) (agency.Message, error) {
						var pp struct{ A, B int }
						if err := json.Unmarshal(params, &pp); err != nil {
							return nil, err
						}
						return agency.NewTextMessage(
							agency.ToolRole,
							fmt.Sprintf("%d", (pp.A+pp.B)*10),
						), nil // *10 is just to distinguish from normal response
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
		agency.NewTextMessage(agency.UserRole, "what is the meaning of life?"),
	)
	if err != nil {
		panic(err)
	}
	printAnswer(answer)

	// test for second function call
	answer, err = t2tOp.Execute(
		ctx,
		agency.NewTextMessage(agency.UserRole, "1+1?"),
	)
	if err != nil {
		panic(err)
	}
	printAnswer(answer)

	// test for both function calls at the same time
	answer, err = t2tOp.Execute(
		ctx,
		agency.NewTextMessage(agency.UserRole, "1+1 and what is the meaning of life?"),
	)
	if err != nil {
		panic(err)
	}
	printAnswer(answer)
}

func printAnswer(message agency.Message) {
	fmt.Printf(
		"Role: %s; Type: %s; Data: %s\n",
		message.Role(),
		message.Kind(),
		string(message.Content()),
	)
}
