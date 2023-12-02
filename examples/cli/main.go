package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

// usage example: go to the repo root and execute
// go run examples/cli/main.go "I love winter" -prompt "You are professional translator, translate everything you see to Russian" -model "gpt-3.5-turbo" -maxTokens=1000
func main() {
	provider := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	temp := flag.Float64("temp", 0.0, "Temperature value")
	maxTokens := flag.Int("max_tokens", 0, "Maximum number of tokens")
	model := flag.String("model", "gpt-3.5-turbo", "Model name")
	prompt := flag.String("prompt", "You are a helpful assistant", "System message")

	flag.Parse()

	if len(os.Args) < 1 {
		fmt.Println("content argument is required")
		return
	}
	content := os.Args[1]

	result, err := provider.
		TextToText(openai.TextToTextParams{
			Model:       *model,
			Temperature: openai.Temperature(float32(*temp)),
			MaxTokens:   *maxTokens,
		}).
		SetPrompt(*prompt).
		Execute(context.Background(), agency.UserMessage(content))

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
