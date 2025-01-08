package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	go_openai "github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

type UserStruct struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
	City string `json:"city"`
}

func main() {
	op := openai.
		New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		TextToText(openai.TextToTextParams{Model: go_openai.GPT4oMini}).
		SetPrompt("You are a poet that writes poetry about the user's data")

	userMsg, err := agency.NewJSONTextMessage(
		agency.UserRole,
		UserStruct{Name: "John", Age: 30, City: "New York"},
	)
	if err != nil {
		panic(err)
	}

	result, err := op.Execute(context.Background(), userMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result.Content()))
}
