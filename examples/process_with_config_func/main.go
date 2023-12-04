package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

// In this example we demonstrate how we can use config func to build a process where 3 step uses the result of the 1 step
func main() {
	provider := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})
	params := openai.TextToTextParams{
		Model:       "gpt-3.5-turbo",
		Temperature: openai.Temperature(0),
		MaxTokens:   100,
	}

	result, _, err := agency.NewProcess(
		agency.ProcessStep{
			Operation: provider.
				TextToText(params).
				SetPrompt("Increase the number by adding 1 to it. Answer only in numbers, without text"),
		},
		agency.ProcessStep{
			Operation: provider.
				TextToText(params).
				SetPrompt("Double the number. Answer only in numbers, without text"),
		},
		agency.ProcessStep{
			ConfigFunc: func(history agency.ProcessHistory, cfg *agency.OperationConfig) error {
				firstStepResult, _ := history.Get(0)                // we ignore error because it's obvious first step exist at the time third executed
				cfg.Prompt = fmt.Sprintf("Add %s", firstStepResult) // we override the prompt with the result of the first step
				return nil
			},
			Operation: provider.TextToText(params), // Note that we don't use SetPrompt because we already set prompt in config func
		},
	).Execute(
		context.Background(),
		agency.UserMessage("5"),
		logStep,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func logStep(in, out agency.Message, cfg agency.OperationConfig, stepIndex uint) {
	fmt.Printf("---\n\nSTEP %d executed\n\nINPUT: %v\n\nCONFIG: %v\n\nOUTPUT: %v\n\n", stepIndex, in, cfg, out)
}
