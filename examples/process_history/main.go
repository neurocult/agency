package main

import (
	"context"
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	provider := openai.New(openai.Params{Key: "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"})
	params := openai.TextToTextParams{
		Model:       "gpt-3.5-turbo",
		Temperature: openai.Temperature(0),
		MaxTokens:   100,
	}

	result, err := agency.NewProcess(
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
				firstStepResult, _ := history.Get(0)
				cfg.Prompt = fmt.Sprintf("Add %s", firstStepResult)
				return nil
			},
			Operation: provider.TextToText(params),
		},
	).Execute(
		context.Background(),
		agency.UserMessage("5"),
		func(in, out agency.Message, cfg *agency.OperationConfig, stepIndex uint) {
			fmt.Printf("---\n\nSTEP %d executed\n\nINPUT: %v\n\nCONFIG: %v\n\nOUTPUT: %v\n\n", stepIndex, in, cfg, out)
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

// InjectHistory allows to pass history between operations by injecting it into the config.
func InjectHistory(history agency.ProcessHistory, cfg *agency.OperationConfig) error {
	cfg.Messages = history.All()
	return nil
}
