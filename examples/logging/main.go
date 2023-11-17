package main

import (
	"context"
	"fmt"

	"github.com/eqtlab/lib/pipeline"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func Logger(input, output core.Message, options ...core.PipeOption) {
	fmt.Printf("in: %v\nout: %v\noptions: %v\n", input, output, options)
}

func main() {
	openAIClient := goopenai.NewClient("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")

	textPipe := openai.TextToText(openAIClient, openai.TextToTextParams{
		Model:       goopenai.GPT3Dot5Turbo,
		Temperature: 0.5,
	})

	_, err := pipeline.New(
		textPipe.WithOptions(core.WithPrompt("explain what that means")),
		textPipe.WithOptions(core.WithPrompt("translate to russian")),
		textPipe.WithOptions(core.WithPrompt("replace all spaces with '_' ")),
	).
		AfterEach(Logger).
		Execute(context.Background(), core.NewUserMessage("Kazakhstan alga!"))

	if err != nil {
		panic(err)
	}
}
