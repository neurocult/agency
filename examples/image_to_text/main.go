package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	imgBytes, err := os.ReadFile("example.png")
	if err != nil {
		panic(err)
	}

	result, err := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		ImageToText(openai.ImageToTextParams{MaxTokens: 300}).
		SetPrompt("describe what you see").
		Execute(
			context.Background(),
			agency.Message{Content: imgBytes},
		)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
