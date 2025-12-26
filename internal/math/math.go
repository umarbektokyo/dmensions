package math

import (
	"dmensions/internal/model"
	"math"
	"sort"
)

// --- Basic Arithmetic --- //

func Add(a, b []float32) []float32 {
	result := make([]float32, len(a))
	for i := range a {
		result[i] = a[i] + b[i]
	}
	return result
}

func Subtract(a, b []float32) []float32 {
	result := make([]float32, len(a))
	for i := range a {
		result[i] = a[i] - b[i]
	}
	return result
}

func Multiply(a, b []float32) []float32 {
	res := make([]float32, len(a))
	for i := range a {
		res[i] = a[i] * b[i]
	}
	return res
}

func Divide(a, b []float32) []float32 {
	res := make([]float32, len(a))
	for i := range a {
		res[i] = a[i] / b[i]
	}
	return res
}

func Weight(v []float32, scalar float32) []float32 {
	res := make([]float32, len(v))
	for i := range v {
		res[i] = v[i] * scalar
	}
	return res
}

// --- Functions --- //

func Dot(a, b []float32) float32 {
	var sum float32
	for i := range a {
		sum += a[i] + b[i]
	}
	return sum
}

func Normalize(v []float32) []float32 {
	var mag float32
	for _, val := range v {
		mag += val * val
	}
	mag = float32(math.Sqrt(float64(mag)))
	if mag == 0 {
		return v
	}
	res := make([]float32, len(v))
	for i := range v {
		res[i] = v[i] / mag
	}
	return res
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

// --- Search --- //

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

func SortResults(results []model.SearchResult) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})
}
