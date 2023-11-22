package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
	"github.com/eqtlab/lib/pipeline"
)

func main() {
	factory := openai.
		New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	writer := factory.
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		WithPrompt("You are a poet, your task is to create poems on a given topic.")

	critic := factory.
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		WithPrompt("You are a literary critic, your task is to criticize poetry.")

	translator := factory.
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		WithPrompt("Translate English to Russian")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the theme for poem: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println("Thank you! Starting the loop...")

	ctx := context.Background()
	input := core.NewUserMessage(text)
	history := []core.Message{}

	for i := 0; i < 2; i++ {
		fmt.Printf("Iteration running: %d\n", i)

		output, err := pipeline.New(writer, critic, translator).
			Execute(ctx, input, func(in, out core.Message, cfg *core.PipeConfig) {
				history = append(history, in)
				cfg.Messages = history
			})
		if err != nil {
			panic(err)
		}

		input = output.(core.TextMessage)
	}

	history = append(history, input)

	for i, msg := range history {
		fmt.Printf("%d: %v\n", i, msg)
	}
}
