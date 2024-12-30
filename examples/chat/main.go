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
		TextToText(openai.TextToTextParams{Model: "gpt-4o-mini"}).
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

		input := agency.NewMessage(agency.UserRole, agency.TextKind, []byte(text))
		answer, err := assistant.SetMessages(messages).Execute(ctx, input)
		if err != nil {
			panic(err)
		}

		fmt.Println("Assistant:", string(answer.Content()))

		messages = append(messages, input, answer)
	}
}
