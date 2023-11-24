// To make this example work make sure you have speech.ogg file in the root of directory.
// You can use text to speech example to generate speech file.
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

	data, err := os.ReadFile("speech.mp3")
	if err != nil {
		panic(err)
	}

	result, err := factory.SpeechToText(openai.SpeechToTextParams{
		Model: goopenai.Whisper1,
	}).Execute(
		context.Background(),
		agency.Message{
			Content: data,
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
