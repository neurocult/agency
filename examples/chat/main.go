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
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)
	history := []core.Message{}
	assistant := openai.
		New(os.Getenv("OPENAI_API_KEY")).
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		WithPrompt("You are helpful assistant.")

	for {
		fmt.Print("User: ")

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred while reading input. Please try again", err)
			return
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
