package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency/core"
	"github.com/neurocult/agency/openai"
	"github.com/neurocult/agency/pipeline"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})
	params := openai.TextToTextParams{Model: "gpt-3.5-turbo"}

	_, err := pipeline.New(
		factory.TextToText(params).SetPrompt("explain what that means"),
		factory.TextToText(params).SetPrompt("translate to russian"),
		factory.TextToText(params).SetPrompt("replace all spaces with '_'"),
	).
		Execute(
			context.Background(),
			core.NewUserMessage("Kazakhstan alga!"),
			Logger,
		)

	if err != nil {
		panic(err)
	}
}

func Logger(input, output core.Message, cfg *core.PipeConfig) {
	fmt.Printf("in: %v\nprompt: %v\nout: %v\n\n", input, cfg.Prompt, output)
}
