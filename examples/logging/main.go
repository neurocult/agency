package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
	"github.com/eqtlab/lib/pipeline"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})
	params := openai.TextToTextParams{Model: "gpt-3.5-turbo"}

	_, err := pipeline.New(
		factory.TextToText(params).WithPrompt("explain what that means"),
		factory.TextToText(params).WithPrompt("translate to russian"),
		factory.TextToText(params).WithPrompt("replace all spaces with '_'"),
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
