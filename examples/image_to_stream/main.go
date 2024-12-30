package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/neurocult/agency"

	providers "github.com/neurocult/agency/providers/openai"
)

func main() {
	imgBytes, err := os.ReadFile("example.png")
	if err != nil {
		panic(err)
	}

	_, err = providers.New(providers.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		TextToStream(providers.TextToStreamParams{
			TextToTextParams: providers.TextToTextParams{MaxTokens: 300, Model: "gpt-4o"},
			StreamHandler: func(delta, total string, isFirst, isLast bool) error {
				fmt.Println(delta)
				return nil
			}}).
		SetPrompt("describe what you see").
		Execute(
			context.Background(),
			agency.NewMessage(agency.UserRole, agency.ImageKind, imgBytes),
		)
	if err != nil {
		panic(err)
	}
}
