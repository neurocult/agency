package openai

import (
	"github.com/sashabaranov/go-openai"
)

type Factory struct {
	client *openai.Client
}

func New(key string) Factory {
	return Factory{
		client: openai.NewClient(key),
	}
}
