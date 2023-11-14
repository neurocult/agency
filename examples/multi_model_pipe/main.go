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

	gpt3 := openai.NewTextToText(openAIClient, goopenai.GPT3Dot5Turbo)
	gpt4 := openai.NewTextToText(openAIClient, goopenai.GPT4TurboPreview)
	whisper := openai.NewSpeechToText(openAIClient, goopenai.Whisper1)

	hear := whisper() // pipe 1

	summarize := gpt3( // pipe 2
		core.WithTemperature(0.5), // TODO: use temperature in openai
		core.WithMessages(core.NewSystemMessage("summarize: ")),
	)

	capitalize := gpt4( // pipe3
		core.WithMessages(core.NewSystemMessage("capitalize: ")),
	)

	msg, err := hear.
		Then(summarize).
		Then(capitalize).
		Execute(ctx, core.NewSpeechMessage(data))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(msg.Bytes()))
}
