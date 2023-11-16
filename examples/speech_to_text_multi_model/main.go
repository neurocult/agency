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

func (s *Saver) Save(_ context.Context, msg core.Message, _ ...core.PipeOption) (core.Message, error) {
	*s = append(*s, msg)
	return msg, nil
}

func main() {
	data, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	openAIClient := goopenai.NewClient("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")

	thinker := openai.TextToText(openAIClient, openai.TextToTextParams{
		Model:       goopenai.GPT3Dot5Turbo,
		Temperature: 0.5,
	})

	// step 1
	hear := openai.SpeechToText(openAIClient, openai.SpeechToTextParams{Model: goopenai.Whisper1})

	// step 2
	summarize := thinker.WithOptions(core.WithPrompt("summarize the text"))

	// step 3
	translate := thinker.WithOptions(core.WithPrompt("translate the text to spanish"))

	saver := Saver{}

	save := saver.Save

	// execute the whole pipeline
	_, err = hear.
		Then(save).
		Then(summarize).
		Then(save).
		Then(translate).
		Then(save).
		Execute(ctx, core.NewSpeechMessage(data))

	if err != nil {
		panic(err)
	}

	for _, msg := range saver {
		fmt.Println(string(msg.Bytes()))
	}
}
