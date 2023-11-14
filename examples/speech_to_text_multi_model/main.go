package main

import (
	"context"
	"fmt"
	"os"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	data, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	openAIClient := goopenai.NewClient("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")

	// step 1
	hear := openai.SpeechToText(openAIClient, core.SpeechToTextParams{Model: goopenai.Whisper1})

	// step 2
	summarize := openai.
		TextToText(openAIClient, core.TextToTextParams{
			Model:       goopenai.GPT3Dot5Turbo,
			Messages:    []core.TextMessage{core.NewSystemMessage("summarize: ")},
			Temperature: 0.5,
		})

	// step 3
	capitalize := openai.
		TextToText(openAIClient, core.TextToTextParams{
			Model:       goopenai.GPT4TurboPreview,
			Messages:    []core.TextMessage{core.NewSystemMessage("capitalize: ")},
			Temperature: 0,
		})

	// execute the whole pipeline
	msg, err := hear.
		Then(summarize).
		Then(capitalize).
		Execute(ctx, core.NewSpeechMessage(data))

	if err != nil {
		panic(err)
	}

	fmt.Println(string(msg.Bytes()))
}
