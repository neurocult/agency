package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"

	_ "github.com/joho/godotenv/autoload"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	factory := openai.New(os.Getenv("OPENAI_API_KEY"))

	msg, err := factory.TextToImage(openai.TextToImageParams{
		Model:     goopenai.CreateImageModelDallE2,
		ImageSize: goopenai.CreateImageSize256x256,
	}).Execute(
		context.Background(),
		core.NewUserMessage("halloween night at a haunted museum."),
	)

	if err != nil {
		panic(err)
	}

	if err := saveToDisk(msg); err != nil {
		panic(err)
	}

	fmt.Println("img saved to example.png")
}

func saveToDisk(msg core.Message) error {
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