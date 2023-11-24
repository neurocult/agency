// To make this example work make sure you have speech.ogg file in the root of directory
package main

import (
	"bytes"
	"context"
	"image/png"
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

	msg, err := agency.NewProcess(
		factory.SpeechToText(openai.SpeechToTextParams{Model: goopenai.Whisper1}),
		factory.TextToImage(openai.TextToImageParams{
			Model:     goopenai.CreateImageModelDallE2,
			ImageSize: goopenai.CreateImageSize256x256,
		}),
	).Execute(context.Background(), agency.Message{Content: data})
	if err != nil {
		panic(err)
	}

	if err := saveImgToDisk(msg); err != nil {
		panic(err)
	}
}

func saveImgToDisk(msg agency.Message) error {
	r := bytes.NewReader(msg.Content)

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
