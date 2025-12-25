package ai

import (
	"bytes"
	"dmensions/internal/model"
	"dmensions/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetEmbedding(text string) ([]float32, error) {
	url := utils.AIURL
	aimodel := utils.AIMODEL
	payload := model.OllamaEmbedRequest{
		Model: aimodel,
		Input: text,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("could not connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status: %s", resp.Status)
	}

	var result model.OllamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned for input: %s", text)
	}

	return result.Embeddings[0], nil
}
