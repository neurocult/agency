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

	var ai = openai.NewPipeFactory(openAIClient)

	hear := ai.SpeechToText()

	summarize := ai.TextToText(
		core.NewSystemMessage("summarize: "),
	)

	capitalize := ai.TextToText(
		core.NewSystemMessage("capitalize: "),
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
