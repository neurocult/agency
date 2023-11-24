// To make this example work make sure you have speech.ogg file in the root of directory
package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/neurocult/agency/pipeline"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency/core"
	"github.com/neurocult/agency/openai"
)

type Saver []core.Message

func (s *Saver) Save(input, output core.Message, _ *core.PipeConfig) {
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
			Model:       goopenai.GPT3Dot5Turbo,
			Temperature: 0.5,
		}).
		SetPrompt("translate to russian")

	// step 3
	uppercase := factory.
		TextToText(openai.TextToTextParams{
			Model:       goopenai.GPT3Dot5Turbo,
			Temperature: 1,
		}).
		SetPrompt("uppercase every letter of the text")

	saver := Saver{}

	sound, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	speechMsg := core.NewSpeechMessage(sound)

	_, err = pipeline.New(
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
