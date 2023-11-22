package main

import (
	"context"
	"fmt"
	"os"

	"github.com/eqtlab/lib/pipeline"
	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
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
		WithPrompt("translate to russian")

	// step 3
	uppercase := factory.
		TextToText(openai.TextToTextParams{
			Model:       goopenai.GPT3Dot5Turbo,
			Temperature: 1,
		}).
		WithPrompt("uppercase every letter of the text")

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
