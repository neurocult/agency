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

	var factory core.PipeFactory = openai.NewPipeFactory(openAIClient)

	systemMsg := core.NewSystemMessage("You are a helpful assistant that translates English to French")

	pipe := factory.TextToText(systemMsg)

	userMsg := core.NewUserMessage("I love programming.")

	resultMsg, err := pipe(context.Background(), userMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resultMsg.Bytes()))
}
