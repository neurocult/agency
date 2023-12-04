// To make this example work make sure you have speech.ogg file in the root of directory
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
	provider := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	// step 1
	hear := provider.
		SpeechToText(openai.SpeechToTextParams{
			Model: goopenai.Whisper1,
		})

	// step2
	translate := provider.
		TextToText(openai.TextToTextParams{
			Model:       "gpt-3.5-turbo",
			Temperature: openai.Temperature(0.5),
		}).
		SetPrompt("translate to russian")

	// step 3
	uppercase := provider.
		TextToText(openai.TextToTextParams{
			Model:       "gpt-3.5-turbo",
			Temperature: openai.Temperature(1),
		}).
		SetPrompt("uppercase every letter of the text")

	sound, err := os.ReadFile("speech.mp3")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	speechMsg := agency.Message{Content: sound}

	_, history, err := agency.ProcessFromOperations(
		hear,
		translate,
		uppercase,
	).Execute(ctx, speechMsg)
	if err != nil {
		panic(err)
	}

	for _, msg := range history.All() {
		fmt.Println(msg.String())
	}
}
