package openai

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/neurocult/agency"
	"github.com/sashabaranov/go-openai"
)

type Embedding []float32

func EmbeddingToBytes(dimensions int, embeddings []Embedding) ([]byte, error) {
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("embeddings is empty")
	}

	buf := make([]byte, len(embeddings)*dimensions*4)

	for i, embedding := range embeddings {
		if len(embedding) != dimensions {
			return nil, fmt.Errorf("invalid embedding length: %d, expected %d", len(embedding), dimensions)
		}

		for j, f := range embedding {
			u := math.Float32bits(f)
			binary.LittleEndian.PutUint32(buf[(i*dimensions+j)*4:], u)
		}
	}

	return buf, nil
}

func BytesToEmbedding(dimensions int, buf []byte) ([]Embedding, error) {
	if mltp := len(buf) % (dimensions * 4); mltp != 0 {
		return nil, fmt.Errorf("invalid buffer length: got %d, but expected multiple of %d", len(buf), dimensions*4)
	}

	embeddings := make([]Embedding, len(buf)/dimensions/4)
	for i := range embeddings {
		embeddings[i] = make([]float32, dimensions)
		for j := 0; j < dimensions; j++ {
			index := (i*dimensions + j) * 4

			if index+4 > len(buf) {
				return nil, fmt.Errorf("buffer is too small for expected number of embeddings")
			}

			embeddings[i][j] = math.Float32frombits(binary.LittleEndian.Uint32(buf[index:]))
		}
	}

	return embeddings, nil
}

// NullableFloat32 is a type that exists to distinguish between undefined values and real zeros.
// It fixes sashabaranov/go-openai issue with zero temp not included in api request due to how json unmarshal work.
type NullableFloat32 *float32

// Temperature is just a tiny helper to create nullable float32 value from regular float32
func Temperature(v float32) NullableFloat32 {
	return &v
}

// nullableToFloat32 replaces nil with zero (in this case value won't be included in api request)
// and for real zeros it returns math.SmallestNonzeroFloat32 that is as close to zero as possible.
func nullableToFloat32(v NullableFloat32) float32 {
	if v == nil {
		return 0
	}
	if *v == 0 {
		return math.SmallestNonzeroFloat32
	}
	return *v
}

// textMessageToOpenAI works with any agency message but ignores everything except role and content.
func textMessageToOpenAI(message agency.Message) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role:    string(message.Role()),
		Content: string(message.Content()),
	}
}

// agencyToOpenAIMessages returns slice of openai chat completion messages created from given config and message.
// Resulting slices starts with config prompt followed by config messages and ends with given message.
func agencyToOpenAIMessages(cfg *agency.OperationConfig, msg agency.Message) ([]openai.ChatCompletionMessage, error) {
	openAIMessages := append(
		make([]openai.ChatCompletionMessage, 0, len(cfg.Messages)+2),
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: cfg.Prompt,
		},
	)

	for _, cfgMsg := range cfg.Messages {
		openAIMessages = append(openAIMessages, textMessageToOpenAI(cfgMsg))
	}

	openAIMessages = append(openAIMessages, openai.ChatCompletionMessage{
		Role:    string(msg.Role()),
		Content: string(msg.Content()),
	})

	return openAIMessages, nil
}
