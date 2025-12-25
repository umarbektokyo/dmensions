package math

import (
	"dmensions/internal/model"
	"math"
)

func Subtract(a, b []float32) []float32 {
	result := make([]float32, len(a))
	for i := range a {
		result[i] = a[i] - b[i]
	}
	return result
}

func CosineSimilarity(a, b []float32) float32 {
	var dotProduct, magASq, magBSq float32
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magASq += a[i] * a[i]
		magBSq += b[i] * b[i]
	}

	magA := float32(math.Sqrt(float64(magASq)))
	magB := float32(math.Sqrt(float64(magBSq)))

	if magA == 0 || magB == 0 {
		return 0
	}

	return dotProduct / (magA * magB)
}

func Search(queryVector []float32, allWords []model.WordData) []model.SearchResult {
	var results []model.SearchResult

	for _, entry := range allWords {
		score := CosineSimilarity(queryVector, entry.Vector)
		results = append(results, model.SearchResult{
			Word:       entry.Word,
			Similarity: score,
		})
	}
	return results
}
