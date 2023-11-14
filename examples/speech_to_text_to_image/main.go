package main

import (
	"bytes"
	"context"
	"image/png"
	"os"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	data, err := os.ReadFile("speech.ogg")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	openAIClient := goopenai.NewClient("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")

	// step 1
	hear := openai.
		SpeechToText(openAIClient, goopenai.Whisper1).
		Pipe()

	// step 2
	draw := openai.
		TextToImage(openAIClient, goopenai.CreateImageModelDallE2).
		Pipe(core.WithSize(goopenai.CreateImageSize256x256))

	// execute the whole pipeline
	msg, err := hear.
		Then(draw).
		Execute(ctx, core.NewSpeechMessage(data))
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

}
