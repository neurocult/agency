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
	openAIClient := goopenai.NewClient("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")

	var factory core.Configurator = openai.NewPipeFactory(openAIClient)

	pipe := factory.SpeechToText()

	data, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	msg, err := pipe(
		context.Background(),
		core.NewSpeechMessage(data),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(msg.Bytes()))
}
