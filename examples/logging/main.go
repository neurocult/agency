package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
	"github.com/eqtlab/lib/pipeline"
)

func Logger(input, output core.Message, options ...core.PipeOption) {
	fmt.Printf("in: %v\nout: %v\noptions: %v\n", input, output, options)
}

func main() {
	factory := openai.New(os.Getenv("OPENAI_API_KEY"))

	textPipe := factory.TextToText(openai.TextToTextParams{
		Model:       goopenai.GPT3Dot5Turbo,
		Temperature: 0.5,
	})

	_, err := pipeline.New(
		textPipe.WithOptions(core.WithPrompt("explain what that means")),
		textPipe.WithOptions(core.WithPrompt("translate to russian")),
		textPipe.WithOptions(core.WithPrompt("replace all spaces with '_'")),
	).
		AfterEach(Logger).
		Execute(context.Background(), core.NewUserMessage("Kazakhstan alga!"))

	if err != nil {
		panic(err)
	}
}
