package main

import (
	"context"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	input := agency.UserMessage(`
		One does not simply walk into Mordor.
		Its black gates are guarded by more than just Orcs.
		There is evil there that does not sleep, and the Great Eye is ever watchful.
	`)

	msg, err := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")}).
		TextToSpeech(openai.TextToSpeechParams{
			Model:          "tts-1",
			ResponseFormat: "mp3",
			Speed:          1,
			Voice:          "alloy",
		}).
		Execute(context.Background(), input)

	if err != nil {
		panic(err)
	}

	if err := saveToDisk(msg); err != nil {
		panic(err)
	}
}

func saveToDisk(msg agency.Message) error {
	file, err := os.Create("speech.mp3")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(msg.Content)
	if err != nil {
		return err
	}

	return nil
}
