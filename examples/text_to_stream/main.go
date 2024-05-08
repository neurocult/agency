package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})
	stream := make(chan string)

	go func() {
		defer close(stream)

		result, err := factory.
			TextToStream(openai.TextToStreamParams{Model: goopenai.GPT3Dot5Turbo, Stream: stream}).
			SetPrompt("Write a few sentences about topic").
			Execute(context.Background(), agency.NewMessage(agency.UserRole, agency.TextKind, []byte("I love programming.")))

		if err != nil {
			panic(err)
		}

		fmt.Println("\nFinal result:", string(result.Content()))
	}()

	for s := range stream {
		fmt.Println(s)
	}
}
