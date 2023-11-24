package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency/core"
	"github.com/neurocult/agency/openai"
	"github.com/neurocult/agency/pipeline"
)

// FIXME this example probably does not work because interceptor is executed after and not before the pipe
func main() {
	factory := openai.
		New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	poet := factory.
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		SetPrompt("You are a poet, your task is to create poems on a given topic.")

	critic := factory.
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		SetPrompt("You are a literary critic, your task is to criticize poetry.")

	translator := factory.
		TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).
		SetPrompt("Translate English to Russian")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the word: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println("Thank you! Starting the loop...")

	ctx := context.Background()
	input := core.NewUserMessage(text)
	history := []core.Message{}
	pipeLine := pipeline.New(poet, critic, translator)

	for i := 0; i < 2; i++ {
		fmt.Printf("Iteration running: %d\n", i)

		output, err := pipeLine.Execute(ctx, input, func(in, out core.Message, cfg *core.PipeConfig) {
			history = append(history, in)
			cfg.Messages = history
		})
		if err != nil {
			panic(err)
		}

		input = output.(core.TextMessage)
	}

	// HACK: on last iteration when we say "input = output" we never append that input to the history
	// so we have to do that after the for loop
	history = append(history, input)

	// IDEA
	// InterceptorParams{ in, out, cfg, pipes, idx }
	// knowing len(pipes) and idx we can determine whether the iteration is last and append not only input but also output

	for i, msg := range history {
		fmt.Printf("%d: %v\n\n", i, msg)
	}
}
