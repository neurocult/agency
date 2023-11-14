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

	var factory core.Config = openai.NewPipeFactory(openAIClient)

	systemMsg := core.NewSystemMessage("You are a helpful assistant that translates %s to %s").Bind("English", "French")

	pipe := factory.TextToText(systemMsg)

	boundUserMsg := core.NewUserMessage("%s").Bind("I love programming.")

	resultMsg, err := pipe(context.Background(), boundUserMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resultMsg.Bytes()))
}
