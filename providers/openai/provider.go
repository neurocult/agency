package openai

import (
	"github.com/sashabaranov/go-openai"
)

type Provider struct {
	client *openai.Client
}

type Params struct {
	Key     string // Required if not using local LLM.
	BaseURL string // Optional. If not set then default openai base url is used
}

func New(params Params) *Provider {
	cfg := openai.DefaultConfig(params.Key)
	if params.BaseURL != "" {
		cfg.BaseURL = params.BaseURL
	}
	return &Provider{
		client: openai.NewClientWithConfig(cfg),
	}
}
