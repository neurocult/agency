package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"strconv"

	"github.com/emil14/ai/lib"
	"github.com/sashabaranov/go-openai"
)

// === USE-CASES ===

// easy:
// continue text
// rewrite text
// templating (dynamic translator e.g.)
// voice -> text
// text -> image
// image -> text
// voice -> text -> image
// image -> text -> voice
// function call (print user input message)

// medium:
// chat (use assistant API?)
// chat with voice from user
// chat with voice from both user and system
// text -> summarize with small gpt3 -> answer with gpt4
// image -> image (edit image / variations)
// sequential multiple function calls
// concurrent multiple function calls
// conditional function calls (depending on user's input)

// hard:
// RAG - (file-system, vector database, relational database)
// fine-tuning - (make it answer in a desired form)
// phone-call - (similar to how chatGPT works)
// group-function calls - (like google meet works, maybe as a google meet member)

// extra-hard (autonomous agents)
// find a job for me
// help me with my coding pet-project
// help me with my book
// be my life-copilot (read my notes, give me insights)

func main() {
	fmt.Println("=== starting test run ===")

	ctx := context.Background()
	openaiClient := openai.NewClient("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")

	testCases := []func() error{
		// pipe demo - increment 3 times
		func() error {
			var pipe lib.Pipe = func(ctx context.Context, msg lib.Message) (lib.Message, error) {
				i, err := strconv.ParseInt(string(msg.Bytes()), 10, 10)
				if err != nil {
					return nil, err
				}
				inc := strconv.Itoa(int(i) + 1)
				return lib.SystemMessage(inc), nil
			}

			pipe = pipe.Chain(pipe).Chain(pipe)

			msg, err := pipe(ctx, lib.UserMessage("0"))
			if err != nil {
				return err
			}
			fmt.Println(string(msg.Bytes()))

			return nil
		},

		// continue text
		func() error {
			pipe := lib.TextPipe(openaiClient)
			msg, err := pipe(ctx, lib.UserMessage("What is the capital of the great Britain?"))
			if err != nil {
				return err
			}
			fmt.Println(string(msg.Bytes()))
			return nil
		},

		// rewrite (translate from english to french) text
		func() error {
			systemMsg := lib.SystemMessage("You are a helpful assistant that translates English to French")

			pipe := lib.TextPipe(openaiClient, systemMsg)
			msg, err := pipe(ctx, lib.UserMessage("I love programming."))
			if err != nil {
				return err
			}
			fmt.Println(string(msg.Bytes()))
			return nil
		},

		// templating
		func() error {
			systemMsg := lib.SystemMessage("You are a helpful assistant that translates %s to %s").Bind("English", "French")
			pipe := lib.TextPipe(openaiClient, systemMsg)

			msg, err := pipe(ctx, lib.UserMessage("%s").Bind("I love programming."))
			if err != nil {
				return err
			}
			fmt.Println(string(msg.Bytes()))

			return nil
		},

		// text2img
		func() error {
			pipe := lib.ImagePipe(openaiClient)

			msg, err := pipe(
				ctx,
				lib.UserMessage("halloween night at a haunted museum."),
			)
			if err != nil {
				return err
			}

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
		},
	}

	for _, tc := range testCases {
		if err := tc(); err != nil {
			panic(err)
		}
	}

	fmt.Println("=== all tests passed ===")
}
