package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"

	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/openai"
)

func main() {
	factory := openai.New(openai.Params{Key: os.Getenv("OPENAI_API_KEY")})

	msg, err := factory.TextToImage(openai.TextToImageParams{
		Model:     goopenai.CreateImageModelDallE2,
		ImageSize: goopenai.CreateImageSize256x256,
	}).Execute(
		context.Background(),
		agency.UserMessage("halloween night at a haunted museum."),
	)

	if err != nil {
		panic(err)
	}

	if err := saveToDisk(msg); err != nil {
		panic(err)
	}

	fmt.Println("img saved to example.png")
}

func saveToDisk(msg agency.Message) error {
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
