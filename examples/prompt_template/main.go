package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	resultMsg, err := factory.
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		WithPrompt(
			"You are a helpful assistant that translates %s to %s",
			"English", "French",
		).
		Execute(
			context.Background(),
			core.NewUserMessage("%s").Bind("I love programming."),
		)

	if err != nil {
		panic(err)
	}

	fmt.Println(resultMsg)
}
