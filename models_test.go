package gigachat

import (
	"encoding/json"
	"testing"
)

func TestGetGenerationModels(t *testing.T) {
	models := GetGenerationModels()
	
	expectedModels := []string{
		GigaChat2,
		GigaChat2Pro,
		GigaChat2Max,
		GigaChat,
		GigaChatPro,
		GigaChatMax,
	}
	
	if len(models) != len(expectedModels) {
		t.Errorf("GetGenerationModels() returned %d models, want %d", len(models), len(expectedModels))
	}
	
	for i, model := range models {
		if model != expectedModels[i] {
			t.Errorf("GetGenerationModels()[%d] = %v, want %v", i, model, expectedModels[i])
		}
	}
}

func TestGetEmbeddingModels(t *testing.T) {
	models := GetEmbeddingModels()
	
	expectedModels := []string{
		Embeddings,
		EmbeddingsGigaR,
	}
	
	if len(models) != len(expectedModels) {
		t.Errorf("GetEmbeddingModels() returned %d models, want %d", len(models), len(expectedModels))
	}
	
	for i, model := range models {
		if model != expectedModels[i] {
			t.Errorf("GetEmbeddingModels()[%d] = %v, want %v", i, model, expectedModels[i])
		}
	}
}

func TestIsValidGenerationModel(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		expected bool
	}{
		{"valid GigaChat2", GigaChat2, true},
		{"valid GigaChat2Pro", GigaChat2Pro, true},
		{"valid GigaChat2Max", GigaChat2Max, true},
		{"valid GigaChat", GigaChat, true},
		{"valid GigaChatPro", GigaChatPro, true},
		{"valid GigaChatMax", GigaChatMax, true},
		{"invalid model", "InvalidModel", false},
		{"empty string", "", false},
		{"embedding model", Embeddings, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidGenerationModel(tt.model); got != tt.expected {
				t.Errorf("IsValidGenerationModel(%v) = %v, want %v", tt.model, got, tt.expected)
			}
		})
	}
}

func TestIsValidEmbeddingModel(t *testing.T) {
	tests := []struct {
		name     string
		model    string
		expected bool
	}{
		{"valid Embeddings", Embeddings, true},
		{"valid EmbeddingsGigaR", EmbeddingsGigaR, true},
		{"invalid model", "InvalidModel", false},
		{"empty string", "", false},
		{"generation model", GigaChat, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmbeddingModel(tt.model); got != tt.expected {
				t.Errorf("IsValidEmbeddingModel(%v) = %v, want %v", tt.model, got, tt.expected)
			}
		})
	}
}

func TestMessage_JSON(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "Hello, world!",
	}
	
	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal Message: %v", err)
	}
	
	var decoded Message
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal Message: %v", err)
	}
	
	if decoded.Role != msg.Role {
		t.Errorf("Decoded Role = %v, want %v", decoded.Role, msg.Role)
	}
	if decoded.Content != msg.Content {
		t.Errorf("Decoded Content = %v, want %v", decoded.Content, msg.Content)
	}
}

func TestChatRequest_JSON(t *testing.T) {
	temp := 0.7
	maxTokens := 100
	
	req := ChatRequest{
		Model: GigaChat,
		Messages: []Message{
			{Role: "user", Content: "test"},
		},
		Temperature: &temp,
		MaxTokens:   &maxTokens,
		Stream:      false,
	}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ChatRequest: %v", err)
	}
	
	var decoded ChatRequest
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ChatRequest: %v", err)
	}
	
	if decoded.Model != req.Model {
		t.Errorf("Decoded Model = %v, want %v", decoded.Model, req.Model)
	}
	if decoded.Temperature == nil || *decoded.Temperature != *req.Temperature {
		t.Errorf("Decoded Temperature = %v, want %v", decoded.Temperature, req.Temperature)
	}
	if decoded.MaxTokens == nil || *decoded.MaxTokens != *req.MaxTokens {
		t.Errorf("Decoded MaxTokens = %v, want %v", decoded.MaxTokens, req.MaxTokens)
	}
	if decoded.Stream != req.Stream {
		t.Errorf("Decoded Stream = %v, want %v", decoded.Stream, req.Stream)
	}
}

func TestChatResponse_JSON(t *testing.T) {
	jsonStr := `{
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "Hello!"
				},
				"finish_reason": "stop"
			}
		],
		"created": 1234567890,
		"model": "GigaChat",
		"usage": {
			"prompt_tokens": 10,
			"completion_tokens": 5,
			"total_tokens": 15
		},
		"object": "chat.completion"
	}`
	
	var resp ChatResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("Failed to unmarshal ChatResponse: %v", err)
	}
	
	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "Hello!" {
		t.Errorf("Expected content 'Hello!', got %v", resp.Choices[0].Message.Content)
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("Expected total tokens 15, got %d", resp.Usage.TotalTokens)
	}
}

func TestModelsResponse_JSON(t *testing.T) {
	jsonStr := `{
		"data": [
			{
				"id": "GigaChat",
				"object": "model",
				"owned_by": "sberbank"
			}
		],
		"object": "list"
	}`
	
	var resp ModelsResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("Failed to unmarshal ModelsResponse: %v", err)
	}
	
	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 model, got %d", len(resp.Data))
	}
	if resp.Data[0].ID != "GigaChat" {
		t.Errorf("Expected model ID 'GigaChat', got %v", resp.Data[0].ID)
	}
}

func TestTokenResponse_JSON(t *testing.T) {
	jsonStr := `{
		"access_token": "test_token_123",
		"expires_at": 1234567890
	}`
	
	var resp TokenResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("Failed to unmarshal TokenResponse: %v", err)
	}
	
	if resp.AccessToken != "test_token_123" {
		t.Errorf("Expected access token 'test_token_123', got %v", resp.AccessToken)
	}
	if resp.ExpiresAt != 1234567890 {
		t.Errorf("Expected expires at 1234567890, got %d", resp.ExpiresAt)
	}
}

func TestImageResult_JSON(t *testing.T) {
	result := ImageResult{
		Content: "base64_encoded_image_data",
		FileID:  "file_123",
		Response: ChatResponse{
			Model:   GigaChat,
			Created: 1234567890,
		},
	}
	
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ImageResult: %v", err)
	}
	
	var decoded ImageResult
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ImageResult: %v", err)
	}
	
	if decoded.Content != result.Content {
		t.Errorf("Decoded Content = %v, want %v", decoded.Content, result.Content)
	}
	if decoded.FileID != result.FileID {
		t.Errorf("Decoded FileID = %v, want %v", decoded.FileID, result.FileID)
	}
}
