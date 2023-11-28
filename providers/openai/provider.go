package openai

import (
	"math"

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

// NullableFloat32 is a type that exists to distinguish between undefined values and real zeros.
// It fixes sashabaranov/go-openai issue with zero temp not included in api request due to how json unmarshal work.
type NullableFloat32 *float32

// Temperature is just a tiny helper to create nullable float32 value from regular float32
func Temperature(v float32) NullableFloat32 {
	return &v
}

// getTemperature replaces nil with zero (in this case value won't be included in api request)
// and for real zeros it returns math.SmallestNonzeroFloat32 that is as close to zero as possible.
func getTemperature(v NullableFloat32) float32 {
	if v == nil {
		return 0
	}
	if *v == 0 {
		return math.SmallestNonzeroFloat32
	}
	return *v
}
