package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	provider := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})
	params := openai.TextToTextParams{
		Model:       "gpt-3.5-turbo",
		Temperature: openai.Temperature(0),
		MaxTokens:   100,
	}

	var (
		msg         agency.Message   = agency.UserMessage("Flying cars") // this is what we start with
		chatHistory []agency.Message                                     // we gonna use this to accumulate history between iterations
	)

	// FIXME currently this example does not work
	// maybe we need to be able to set initial history for the process to implement this easily
	// current approach cannot be working because on iteration 2 steps 2 and 3 are missing the history of the previous iterations
	// also could be related to the fact that history lacks initial message
	for i := 0; i < 3; i++ {
		output, curHistory, err := agency.NewProcess( // on each iteration we create new process with its own execution context
			// writer
			agency.ProcessStep{
				ConfigFunc: injectHistory,
				Operation: provider.
					TextToText(params).
					SetPrompt("Create slogan for a given context. A short but catchy phrase").
					SetMessages(chatHistory), // start with the history from the previous iteration
			},
			// critic
			agency.ProcessStep{
				ConfigFunc: injectHistory, // use history from the previous iteration plus previous step
				Operation: provider.
					TextToText(params).
					SetPrompt("Criticize the given slogan. Find its weaknesses and suggest improvements"),
			},
			// censor
			agency.ProcessStep{
				ConfigFunc: injectHistory, // use history from the previous iteration plus two previous steps
				Operation: provider.
					TextToText(params).
					SetPrompt(
						"You are a safe guard. Text must not contain expectations about future. If you see anything happy, point to it so it can be removed.",
					),
			},
		).Execute(
			context.Background(),
			msg,
			logStep,
		)

		if err != nil {
			panic(err)
		}

		chatHistory = curHistory.All()
		chatHistory = chatHistory[0 : len(chatHistory)-1] // remove last output to avoid duplication, we will have it as the input
		msg = output
	}

	fmt.Printf("RESULT: %v\n\n", msg)
}

// injectHistory uses history passed in by the process
// and injects in into the configuration of the operation so operation handler has access to history.
func injectHistory(history agency.ProcessHistory, cfg *agency.OperationConfig) error {
	cfg.Messages = append(cfg.Messages, history.All()...)
	return nil
}

// logStep simply prints data related to each step execution. It implements interceptor interface.
func logStep(in, out agency.Message, cfg agency.OperationConfig, stepIndex uint) {
	fmt.Printf("---\n\nSTEP %d executed\n\nINPUT: %v\n\nCONFIG: %v\n\nOUTPUT: %v\n\n", stepIndex, in, cfg, out)
}
