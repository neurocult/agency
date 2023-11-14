package main

import (
	"context"
	"fmt"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	openAIClient := goopenai.NewClient("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")

	var factory core.Configurator = openai.NewPipeFactory(openAIClient)

	pipe := factory.TextToText()
	userMsg := core.NewUserMessage("What is the capital of the great Britain?")

	resultMsg, err := pipe(context.Background(), userMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resultMsg.Bytes()))
}
