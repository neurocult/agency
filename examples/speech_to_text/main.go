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

	data, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	result, err := factory.SpeechToText(openai.SpeechToTextParams{
		Model: goopenai.Whisper1,
	}).Execute(
		context.Background(),
		core.NewSpeechMessage(data),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
