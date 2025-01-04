package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	result, err := factory.
		TextToStream(openai.TextToStreamParams{
			TextToTextParams: openai.TextToTextParams{Model: goopenai.GPT4oMini},
			StreamHandler: func(delta, total string, isFirst, isLast bool) error {
				if isFirst {
					fmt.Println("====Start streaming====")
				}
				fmt.Print(delta)
				if isLast {
					fmt.Println("\n====Finish streaming====")
				}
				return nil
			},
		}).
		SetPrompt("Write a few sentences about topic").
		Execute(
			context.Background(),
			agency.NewTextMessage(agency.UserRole, "I love programming."),
		)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nResult:", string(result.Content()))
}
