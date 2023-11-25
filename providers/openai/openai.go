package openai

import (
	"github.com/sashabaranov/go-openai"
)

type Factory struct {
	client *openai.Client
}

type Params struct {
	Key     string // Required if not using local LLM.
	BaseURL string // Optional. If not set then default openai base url is used
}

func New(params Params) *Factory {
	cfg := openai.DefaultConfig(params.Key)
	if params.BaseURL != "" {
		cfg.BaseURL = params.BaseURL
	}
	return &Factory{
		client: openai.NewClientWithConfig(cfg),
	}
}
