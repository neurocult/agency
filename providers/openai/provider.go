package openai

import (
	"github.com/sashabaranov/go-openai"
)

// Provider is a set of operation builders.
type Provider struct {
	client *openai.Client
}

// Params is a set of parameters specific for creating this concrete provider.
// They are shared across all operation builders.
type Params struct {
	Key     string // Required if not using local LLM.
	BaseURL string // Optional. If not set then default openai base url is used
}

// New creates a new Provider instance.
func New(params Params) *Provider {
	cfg := openai.DefaultConfig(params.Key)
	if params.BaseURL != "" {
		cfg.BaseURL = params.BaseURL
	}
	return &Provider{
		client: openai.NewClientWithConfig(cfg),
	}
}
