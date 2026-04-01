package minimax

import (
	"testing"

	agencyopenai "github.com/neurocult/agency/providers/openai"
)

func TestNew(t *testing.T) {
	p := New(Params{Key: "test-key"})
	if p == nil {
		t.Fatal("New returned nil")
	}
	if p.inner == nil {
		t.Fatal("New returned provider with nil inner client")
	}
}

func TestTextToTextReturnsOperation(t *testing.T) {
	p := New(Params{Key: "test-key"})
	op := p.TextToText(TextToTextParams{Model: ModelM27})
	if op == nil {
		t.Fatal("TextToText returned nil operation")
	}
}

func TestTextToStreamReturnsOperation(t *testing.T) {
	p := New(Params{Key: "test-key"})
	op := p.TextToStream(TextToStreamParams{
		TextToTextParams: TextToTextParams{Model: ModelM27HighSpeed},
	})
	if op == nil {
		t.Fatal("TextToStream returned nil operation")
	}
}

func TestClampTemp(t *testing.T) {
	tests := []struct {
		name  string
		input agencyopenai.NullableFloat32
		want  float32
	}{
		{
			name:  "nil is clamped to minTemp",
			input: nil,
			want:  minTemp,
		},
		{
			name:  "zero is clamped to minTemp",
			input: agencyopenai.Temperature(0),
			want:  minTemp,
		},
		{
			name:  "negative is clamped to minTemp",
			input: agencyopenai.Temperature(-0.5),
			want:  minTemp,
		},
		{
			name:  "valid low value is unchanged",
			input: agencyopenai.Temperature(0.1),
			want:  0.1,
		},
		{
			name:  "valid mid value is unchanged",
			input: agencyopenai.Temperature(0.7),
			want:  0.7,
		},
		{
			name:  "exactly 1.0 is valid and unchanged",
			input: agencyopenai.Temperature(1.0),
			want:  1.0,
		},
		{
			name:  "above 1.0 is clamped to 1.0",
			input: agencyopenai.Temperature(1.5),
			want:  1.0,
		},
		{
			name:  "well above 1.0 is clamped to 1.0",
			input: agencyopenai.Temperature(2.0),
			want:  1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clampTemp(tt.input)
			if got == nil {
				t.Fatal("clampTemp returned nil")
			}
			if *got != tt.want {
				t.Errorf("clampTemp(%v) = %v, want %v", tt.input, *got, tt.want)
			}
		})
	}
}

func TestModelConstants(t *testing.T) {
	models := []string{ModelM27, ModelM27HighSpeed, ModelM25, ModelM25HighSpeed}
	for _, m := range models {
		if m == "" {
			t.Errorf("model constant is empty")
		}
	}
}

func TestBaseURL(t *testing.T) {
	if BaseURL != "https://api.minimax.io/v1" {
		t.Errorf("unexpected BaseURL: %s", BaseURL)
	}
}
