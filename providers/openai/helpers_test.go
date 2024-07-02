package openai

import (
	"reflect"
	"testing"
)

func TestEmbeddingToBytes(t *testing.T) {
	floats := []Embedding{{1.1, 2.2, 3.3}, {4.4, 5.5, 6.6}}

	bytes, err := EmbeddingToBytes(3, floats)
	if err != nil {
		t.Errorf("EmbeddingToBytes error %v", err)
	}

	newFloats, err := BytesToEmbedding(3, bytes)
	if err != nil {
		t.Errorf("EmbeddingToBytes error %v", err)
	}

	if !reflect.DeepEqual(floats, newFloats) {
		t.Errorf("floats and newFloats are not equal %v %v", floats, newFloats)
	}

	wrongFloats := []Embedding{{4.4, 5.5, 6.6, 7.7}}

	_, err = EmbeddingToBytes(3, wrongFloats)
	if err == nil {
		t.Errorf("EmbeddingToBytes should has error")
	}

	bytes, err = EmbeddingToBytes(4, wrongFloats)
	if err != nil {
		t.Errorf("EmbeddingToBytes error %v", err)
	}

	_, err = BytesToEmbedding(3, bytes)
	if err == nil {
		t.Errorf("BytesToEmbedding should has error")
	}
}
