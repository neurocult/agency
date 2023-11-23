package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	result, err := factory.
		TextToText(openai.TextToTextParams{Model: goopenai.GPT3Dot5Turbo}).
		SetPrompt("You are a helpful assistant that translates English to French").
		Execute(context.Background(), core.NewUserMessage("I love programming."))

	if err != nil {
		panic(err)
	}

	fmt.Println(string(result.Bytes()))
}
