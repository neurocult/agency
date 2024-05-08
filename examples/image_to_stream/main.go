package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/neurocult/agency"

	providers "github.com/neurocult/agency/providers/openai"
	"github.com/sashabaranov/go-openai"
)

func main() {
	imgBytes, err := os.ReadFile("assets/dracula.png")
	if err != nil {
		panic(err)
	}

	stream := make(chan string)

	go func() {
		defer close(stream)
		result, err := providers.New(providers.Params{Key: os.Getenv("OPENAI_API_KEY")}).
			TextToStream(providers.TextToStreamParams{MaxTokens: 300, Model: openai.GPT4Turbo, Stream: stream}).
			SetPrompt("describe what you see").
			Execute(
				context.Background(),
				agency.NewMessage(agency.UserRole, agency.ImageKind, imgBytes),
			)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(result.Content()))
	}()

	for s := range stream {
		fmt.Println(s)
	}
}
