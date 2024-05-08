package openai

import (
	"encoding/binary"
	"fmt"
	"math"
)

func EmbeddingToBytes(dimensions int, embeddings [][]float32) ([]byte, error) {
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

func BytesToEmbedding(dimensions int, buf []byte) ([][]float32, error) {
	if mltp := len(buf) % (dimensions * 4); mltp != 0 {
		return nil, fmt.Errorf("invalid buffer length: got %d, but expected multiple of %d", len(buf), dimensions*4)
	}

	embeddings := make([][]float32, len(buf)/dimensions/4)
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
