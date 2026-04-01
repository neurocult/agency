// Package minimax provides an Agency provider for the MiniMax LLM API.
//
// MiniMax exposes an OpenAI-compatible REST interface, so this package is a
// thin wrapper around the openai provider that:
//   - points the client at the MiniMax base URL (https://api.minimax.io/v1),
//   - clamps temperature into the range (0, 1] required by MiniMax, and
//   - exports convenience constants for the available model names.
//
// Usage:
//
//	p := minimax.New(minimax.Params{Key: os.Getenv("MINIMAX_API_KEY")})
//	op := p.TextToText(minimax.TextToTextParams{Model: minimax.ModelM27}).
//	       SetPrompt("You are a helpful assistant.")
package minimax

import (
	"github.com/neurocult/agency"
	agencyopenai "github.com/neurocult/agency/providers/openai"
)

// BaseURL is the MiniMax OpenAI-compatible API endpoint.
const BaseURL = "https://api.minimax.io/v1"

// Model name constants for available MiniMax models.
// All models support a 204 K context window.
const (
	ModelM27         = "MiniMax-M2.7"
	ModelM27HighSpeed = "MiniMax-M2.7-highspeed"
	ModelM25         = "MiniMax-M2.5"
	ModelM25HighSpeed = "MiniMax-M2.5-highspeed"
)

// minTemp is the smallest temperature value accepted by MiniMax (must be > 0).
const minTemp float32 = 0.01

// Type aliases so callers only need to import this package.
type (
	TextToTextParams   = agencyopenai.TextToTextParams
	TextToStreamParams = agencyopenai.TextToStreamParams
)

// Params contains the configuration required to create a Provider.
type Params struct {
	// Key is the MiniMax API key (MINIMAX_API_KEY).
	Key string
}

// Provider wraps the openai provider and configures it for the MiniMax API.
type Provider struct {
	inner *agencyopenai.Provider
}

// New creates a new MiniMax Provider.
func New(params Params) *Provider {
	return &Provider{
		inner: agencyopenai.New(agencyopenai.Params{
			Key:     params.Key,
			BaseURL: BaseURL,
		}),
	}
}

// TextToText returns an operation that sends a text prompt to a MiniMax model
// and returns the model's text response.
func (p *Provider) TextToText(params TextToTextParams) *agency.Operation {
	params.Temperature = clampTemp(params.Temperature)
	return p.inner.TextToText(params)
}

// TextToStream returns an operation that streams the response from a MiniMax
// model, calling params.StreamHandler for each incremental delta.
func (p *Provider) TextToStream(params TextToStreamParams) *agency.Operation {
	params.TextToTextParams.Temperature = clampTemp(params.TextToTextParams.Temperature)
	return p.inner.TextToStream(params)
}

// clampTemp ensures that the temperature value satisfies MiniMax's constraint
// that temperature must be in the range (0, 1].
// A nil or non-positive value is clamped to minTemp; a value above 1 is
// clamped to 1.
func clampTemp(t agencyopenai.NullableFloat32) agencyopenai.NullableFloat32 {
	if t == nil || *t <= 0 {
		return agencyopenai.Temperature(minTemp)
	}
	if *t > 1 {
		return agencyopenai.Temperature(1)
	}
	return t
}
