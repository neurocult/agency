package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	go_openai "github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	result, err := factory.
		TextToText(openai.TextToTextParams{Model: go_openai.GPT4oMini}).
		SetPrompt("You are a helpful assistant that translates English to French").
		Execute(
			context.Background(),
			agency.NewTextMessage(agency.UserRole, "I love programming."),
		)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(result.Content()))
}
