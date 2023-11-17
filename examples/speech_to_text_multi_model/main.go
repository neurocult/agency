package main

import (
	"context"
	"fmt"
	"os"

	"github.com/eqtlab/lib/pipeline"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

type Saver []core.Message

func (s *Saver) Save(input, output core.Message, _ ...core.PipeOption) {
	*s = append(*s, output)
}

func main() {
	openAIClient := goopenai.NewClient("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")

	// step 1
	hear := openai.
		SpeechToText(openAIClient, openai.SpeechToTextParams{
			Model: goopenai.Whisper1,
		})

	// step2
	translate := openai.
		TextToText(openAIClient, openai.TextToTextParams{
			Model:       goopenai.GPT3Dot5Turbo,
			Temperature: 0.5,
		}).
		WithOptions(core.WithPrompt("translate to russian"))

	// step 3
	uppercase := openai.
		TextToText(openAIClient, openai.TextToTextParams{
			Model:       goopenai.GPT3Dot5Turbo,
			Temperature: 1,
		}).
		WithOptions(core.WithPrompt("uppercase every letter of the text"))

	saver := Saver{}

	sound, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	speechMsg := core.NewSpeechMessage(sound)

	pipeline.New(
		hear,
		translate,
		uppercase,
	).
		InterceptEach(saver.Save).
		Execute(ctx, speechMsg)

	if err != nil {
		panic(err)
	}

	for _, msg := range saver {
		fmt.Println(msg.String())
	}
}
