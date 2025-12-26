package math

import (
	"dmensions/internal/model"
	"math"
	"math/rand"
)

// I must re-inmplement it in a separate library.
// I haven't fully comprehended the algorithms so here's what me and gemini came up with.

// ProjectTSNE runs a simplified t-SNE to reduce dimensions to 2D.
// iter: number of iterations (try 300-500).
// learningRate: speed of movement (try 10.0 - 100.0).
func ProjectTSNE(data []model.WordData, iter int, learningRate float64) map[int64]model.Point2D {
	n := len(data)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return map[int64]model.Point2D{data[0].ID: {X: 0.5, Y: 0.5}}
	}

	// 1. Initialize random 2D positions (sigma 0.0001)
	points := make([]model.Point2D, n)
	for i := range points {
		points[i] = model.Point2D{
			X: rand.Float64() * 0.0001,
			Y: rand.Float64() * 0.0001,
		}
	}

	// 2. Compute Pairwise Affinities in High-D (P matrix)
	p := make([][]float64, n)
	for i := range p {
		p[i] = make([]float64, n)
	}

	const twoSigmaSq = 2.0 * 0.8 // Adjust variance if needed

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				distSq := euclideanDistSq(data[i].Vector, data[j].Vector)
				p[i][j] = math.Exp(-float64(distSq) / twoSigmaSq)
			}
		}
	}

	// Symmetrize P and normalize
	var sumP float64
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			p[i][j] = (p[i][j] + p[j][i]) / (2 * float64(n))
			sumP += p[i][j]
		}
	}
	// Safety check to avoid division by zero
	if sumP == 0 {
		sumP = 1
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			p[i][j] /= sumP
		}
	}

	// 3. Gradient Descent Loop
	for k := 0; k < iter; k++ {
		// A. Compute Low-D Affinities (Q matrix) - Student-t distribution (1 degree of freedom)
		// q[i][j] = (1 + ||yi - yj||^2)^-1
		q := make([][]float64, n)
		for i := range q {
			q[i] = make([]float64, n)
		}

		num := make([][]float64, n) // Numerator part
		for i := range num {
			num[i] = make([]float64, n)
		}

		var sumQ float64
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j {
					distSq := distSq2D(points[i], points[j])
					val := 1.0 / (1.0 + distSq)
					num[i][j] = val
					sumQ += val
				}
			}
		}

		// Normalize Q
		if sumQ == 0 {
			sumQ = 1
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				q[i][j] = num[i][j] / sumQ
			}
		}

		// B. Compute Gradients
		// dC/dyi = 4 * sum( (pij - qij) * num[i][j] * (yi - yj) )
		grads := make([]model.Point2D, n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j {
					mult := 4.0 * (p[i][j] - q[i][j]) * num[i][j]
					grads[i].X += mult * (points[i].X - points[j].X)
					grads[i].Y += mult * (points[i].Y - points[j].Y)
				}
			}
		}

		// C. Update Positions
		for i := 0; i < n; i++ {
			points[i].X -= learningRate * grads[i].X
			points[i].Y -= learningRate * grads[i].Y
		}
	}

	// 4. Map back to ID
	results := make(map[int64]model.Point2D)
	for i, d := range data {
		results[d.ID] = points[i]
	}
	return results
}

// Helpers

func euclideanDistSq(a, b []float32) float32 {
	var sum float32
	// Assumes equal length
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}

func distSq2D(a, b model.Point2D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}
