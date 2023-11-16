package main

import (
	"context"
	"fmt"
	"os"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

type Saver []core.Message

func (s *Saver) Save(msg core.Message, _ ...core.PipeOption) {
	*s = append(*s, msg)
}

func main() {
	openAIClient := goopenai.NewClient("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")

	// step 1
	hear := openai.
		SpeechToText(openAIClient, openai.SpeechToTextParams{
			Model: goopenai.Whisper1,
		})

	// step2
	summarize := openai.
		TextToText(openAIClient, openai.TextToTextParams{
			Model:       goopenai.GPT3Dot5Turbo,
			Temperature: 0.5,
		}).
		WithOptions(core.WithPrompt("summarize the text"))

	// step 3
	capitalize := openai.
		TextToText(openAIClient, openai.TextToTextParams{
			Model:       goopenai.GPT4TurboPreview,
			Temperature: 1,
		}).
		WithOptions(core.WithPrompt("capitalize the text"))

	saver := Saver{}

	sound, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	speechMsg := core.NewSpeechMessage(sound)

	msg, err := hear.
		Intercept(saver.Save).
		Then(summarize).
		Then(capitalize).
		Execute(ctx, speechMsg)

	if err != nil {
		panic(err)
	}

	fmt.Println(saver)

	for _, msg := range saver {
		fmt.Println(string(msg.Bytes()))
	}

	fmt.Println(msg)
}
