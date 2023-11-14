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

	gpt3 := openai.TextToText(openAIClient, goopenai.GPT3Dot5Turbo)

	pipe := gpt3()

	userMsg := core.NewUserMessage("What is the capital of Great Britain?")

	resultMsg, err := pipe(context.Background(), userMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resultMsg.Bytes()))
}
