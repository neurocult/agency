package main

import (
	"bytes"
	"context"
	"image/png"
	"os"

	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
	"github.com/eqtlab/lib/pipeline"
)

func main() {
	factory := openai.New(os.Getenv("OPENAI_API_KEY"))

	data, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	msg, err := pipeline.New(
		factory.SpeechToText(openai.SpeechToTextParams{Model: goopenai.Whisper1}),
		factory.TextToImage(openai.TextToImageParams{
			Model:     goopenai.CreateImageModelDallE2,
			ImageSize: goopenai.CreateImageSize256x256,
		}),
	).Execute(context.Background(), core.NewSpeechMessage(data))
	if err != nil {
		panic(err)
	}

	if err := saveImgToDisk(msg); err != nil {
		panic(err)
	}
}

func saveImgToDisk(msg core.Message) error {
	r := bytes.NewReader(msg.Bytes())

	imgData, err := png.Decode(r)
	if err != nil {
		return err
	}

	file, err := os.Create("example.png")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		return err
	}

	return nil
}
