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

type Saver []agency.Message

// This is how we can retrieve process history by hand with the interceptor, without using the history itself.
// But we can't (or it's hard to do) pass history between steps this way. For that we can use config func.
func (s *Saver) Save(input, output agency.Message, _ *agency.OperationConfig, _ uint) {
	*s = append(*s, output)
}

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

	saver := Saver{}

	sound, err := os.ReadFile("speech.mp3")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	speechMsg := agency.Message{Content: sound}

	_, err = agency.ProcessFromOperations(
		hear,
		translate,
		uppercase,
	).Execute(ctx, speechMsg, saver.Save)
	if err != nil {
		panic(err)
	}

	for _, msg := range saver {
		fmt.Println(msg.String())
	}
}
