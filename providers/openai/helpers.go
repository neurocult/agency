package openai

import (
	"encoding/binary"
	"fmt"
	"math"
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
