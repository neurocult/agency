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

	gpt3 := openai.TextToText(openAIClient, goopenai.GPT3Dot5Turbo)
	gpt4 := openai.TextToText(openAIClient, goopenai.GPT4TurboPreview)
	whisper := openai.SpeechToText(openAIClient, goopenai.Whisper1)

	// pipe 1
	hear := whisper()

	// pipe 2
	summarize := gpt3(
		core.WithTemperature[core.TextConfig](0.5),
		core.WithSystemMessage("summarize: "),
	)

	// pipe3
	capitalize := gpt4(
		core.WithSystemMessage("capitalize: "),
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
