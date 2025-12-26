package model

// --- Database --- //

type WordData struct {
	ID     int64     `json:"id"`
	Word   string    `json:"word"`
	Vector []float32 `json:"vector"`
}

type SearchResult struct {
	Word       string  `json:"word"`
	Similarity float32 `json:"similarity"`
}

// --- Ollama --- //

type OllamaEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type OllamaEmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

// --- Visualisation --- //

type Point2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
