// Example of using the MiniMax provider with the Agency library.
//
// Set the MINIMAX_API_KEY environment variable (or add it to a .env file)
// before running:
//
//	MINIMAX_API_KEY=<your-key> go run examples/minimax/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/minimax"
)

func main() {
	assistant := minimax.
		New(minimax.Params{Key: os.Getenv("MINIMAX_API_KEY")}).
		TextToText(minimax.TextToTextParams{Model: minimax.ModelM27}).
		SetPrompt("You are a helpful assistant.")

	messages := []agency.Message{}
	reader := bufio.NewReader(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Print("User: ")

		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		input := agency.NewTextMessage(agency.UserRole, text)
		answer, err := assistant.SetMessages(messages).Execute(ctx, input)
		if err != nil {
			panic(err)
		}

		fmt.Println("Assistant:", string(answer.Content()))

		messages = append(messages, input, answer)
	}
}
