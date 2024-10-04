package embeddingfunction

import (
	"fmt"
)

type PreEmbedFunction struct {
	EmbedMap map[string][]float32
}

func (e *PreEmbedFunction) CreateEmbedding(documents []string) ([][]float32, error) {
	out := make([][]float32, 0)
	for _, doc := range documents {
		if e.EmbedMap[doc] == nil {
			return nil, fmt.Errorf("doc embedding must be sent over first")
		}
		out = append(out, e.EmbedMap[doc])
	}
	return out, nil
}

func (e *PreEmbedFunction) CreateEmbeddingWithModel(documents []string, model string) ([][]float32, error) {
	out := make([][]float32, 0)
	for _, doc := range documents {
		if e.EmbedMap[doc] == nil {
			return nil, fmt.Errorf("doc embedding must be sent over first")
		}
		out = append(out, e.EmbedMap[doc])
	}
	return out, nil
}
