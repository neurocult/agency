package minimax_test

import (
	"context"
	"os"
	"testing"

	"github.com/neurocult/agency"
	"github.com/neurocult/agency/providers/minimax"
)

func TestIntegrationTextToText(t *testing.T) {
	key := os.Getenv("MINIMAX_API_KEY")
	if key == "" {
		t.Skip("MINIMAX_API_KEY not set")
	}

	p := minimax.New(minimax.Params{Key: key})
	op := p.TextToText(minimax.TextToTextParams{
		Model: minimax.ModelM27HighSpeed,
	}).SetPrompt("You are a helpful assistant. Be very brief.")

	ctx := context.Background()
	input := agency.NewTextMessage(agency.UserRole, "Reply with exactly the word: PONG")
	answer, err := op.Execute(ctx, input)
	if err != nil {
		t.Fatalf("TextToText failed: %v", err)
	}
	if len(answer.Content()) == 0 {
		t.Fatal("got empty response")
	}
	t.Logf("Response: %s", answer.Content())
}
