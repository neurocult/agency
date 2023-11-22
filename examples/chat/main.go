package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	assistant := openai.
		New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		WithPrompt("You are helpful assistant.")

	history := []core.Message{}
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Print("User: ")

		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input := core.NewUserMessage(text)
		answer, err := assistant.WithMessages(history).Execute(ctx, input)
		if err != nil {
			panic(err)
		}

		fmt.Println("Assistant: ", answer)

		history = append(history, input, answer)
	}
}
