package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	openAIClient := goopenai.NewClient("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")

	whisper := openai.TextToImage(openAIClient, goopenai.CreateImageModelDallE2)

	pipe := whisper(
		core.WithSize(goopenai.CreateImageSize256x256),
	)

	msg, err := pipe(
		context.Background(),
		core.NewUserMessage("halloween night at a haunted museum."),
	)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(msg.Bytes())
	imgData, err := png.Decode(r)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("example.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		panic(err)
	}

	fmt.Println("img saved to example.png")
}
