// based on user input
// 1) call 2 functions, A then B then answer
// 2) call 1 function A or B
// 3) just answer (no function call)

package main

import (
	"bufio"
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
					Name:        "SumNumbers",
					Description: "Sum given numbers when asked",
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
		SetPrompt("You are helpful assistant.")

	messages := []agency.Message{}
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Print("User: ")

		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input := agency.UserMessage(text)
		answer, err := t2tOp.SetMessages(messages).Execute(ctx, input)
		if err != nil {
			panic(err)
		}

		fmt.Println("Assistant: ", answer)

		messages = append(messages, input, answer)
	}
}
