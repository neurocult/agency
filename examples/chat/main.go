package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/eqtlab/lib/core"
	"github.com/eqtlab/lib/openai"
)

func main() {
	fmt.Print("Enter text: ")
	reader := bufio.NewReader(os.Stdin)
	systemMesage := "You are helpful assistant."
	assistant := openai.New(os.Getenv("OPENAI_API_KEY")).TextToText(openai.TextToTextParams{Model: "gpt-3.5-turbo"}).WithPrompt(systemMesage)
	history := []core.Message{}

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred while reading input. Please try again", err)
			return
		}

		input := core.NewUserMessage(text)

		answer, err := assistant.
			WithMessages(history).
			Execute(context.Background(), input)
		if err != nil {
			panic(err)
		}

		fmt.Println("Assistant: ", answer)

		history = append(history, input, answer)
	}
}

func saveToDisk(msg core.Message) error {
	file, err := os.Create("example.mp3")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(msg.Bytes())
	if err != nil {
		return err
	}

	return nil
}
