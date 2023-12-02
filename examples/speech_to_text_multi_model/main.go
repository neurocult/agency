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

func (s *Saver) Save(input, output agency.Message, _ *agency.OperationConfig) {
	*s = append(*s, output)
}

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	// step 1
	hear := factory.
		SpeechToText(openai.SpeechToTextParams{
			Model: goopenai.Whisper1,
		})

	// step2
	translate := factory.
		TextToText(openai.TextToTextParams{
			Model:       "gpt-3.5-turbo",
			Temperature: openai.Temperature(0.5),
		}).
		SetPrompt("translate to russian")

	// step 3
	uppercase := factory.
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

	_, err = agency.NewProcess(
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
