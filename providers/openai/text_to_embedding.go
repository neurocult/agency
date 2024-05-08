package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/neurocult/agency"
)

type EmbeddingModel = openai.EmbeddingModel

const AdaEmbeddingV2 EmbeddingModel = openai.AdaEmbeddingV2

type TextToEmbeddingParams struct {
	Model EmbeddingModel
}

func (p Provider) TextToEmbedding(params TextToEmbeddingParams) *agency.Operation {
	return agency.NewOperation(func(ctx context.Context, msg agency.Message, cfg *agency.OperationConfig) (agency.Message, error) {
		//TODO: we have to convert string to model and then model to string. Can we optimize it?
		messages := append(cfg.Messages, msg)
		texts := make([]string, len(messages))

		for i, m := range messages {
			texts[i] = string(m.Content())
		}

		resp, err := p.client.CreateEmbeddings(
			ctx,
			openai.EmbeddingRequest{
				Input: texts,
				Model: params.Model,
			},
		)
		if err != nil {
			return nil, err
		}

		vectors := make([][]float32, len(resp.Data))
		for i, vector := range resp.Data {
			vectors[i] = vector.Embedding
		}

		bytes, err := EmbeddingToBytes(1536, vectors)
		if err != nil {
			return nil, fmt.Errorf("failed to convert embedding to bytes: %w", err)
		}

		//TODO: we have to convert []float32 to []byte. Can we optimize it?
		return agency.NewMessage(agency.AssistantRole, agency.VectorKind, bytes), nil
	})
}