package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	assistant := openai.
		New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
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
		answer, err := assistant.SetMessages(messages).Execute(ctx, input)
		if err != nil {
			panic(err)
		}

		fmt.Println("Assistant: ", answer)

		messages = append(messages, input, answer)
	}
}
