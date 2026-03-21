package gigachat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func createMockTokenManager(token string) *TokenManager {
	tm := NewTokenManager("test_key")
	tm.accessToken = token
	tm.expiresAt = time.Now().Add(1 * time.Hour)
	return tm
}

func TestNewClient(t *testing.T) {
	tm := createMockTokenManager("test_token")
	client := NewClient(tm)
	
	if client.tokenManager != tm {
		t.Error("Expected token manager to be set")
	}
	if client.baseURI != "https://gigachat.devices.sberbank.ru" {
		t.Errorf("Expected default baseURI, got %v", client.baseURI)
	}
	if client.defaultModel != GigaChat {
		t.Errorf("Expected default model %v, got %v", GigaChat, client.defaultModel)
	}
	if client.httpClient == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	tm := createMockTokenManager("test_token")
	customURI := "https://custom.uri"
	customModel := GigaChatPro
	
	client := NewClient(
		tm,
		WithBaseURI(customURI),
		WithDefaultModel(customModel),
		WithClientInsecureSkipVerify(true),
	)
	
	if client.baseURI != customURI {
		t.Errorf("Expected baseURI %v, got %v", customURI, client.baseURI)
	}
	if client.defaultModel != customModel {
		t.Errorf("Expected model %v, got %v", customModel, client.defaultModel)
	}
	
	if transport, ok := client.httpClient.Transport.(*http.Transport); ok {
		if !transport.TLSClientConfig.InsecureSkipVerify {
			t.Error("Expected InsecureSkipVerify to be true")
		}
	}
}

func TestWithBaseURI_TrailingSlash(t *testing.T) {
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI("https://example.com/"))
	
	if client.baseURI != "https://example.com" {
		t.Errorf("Expected trailing slash to be removed, got %v", client.baseURI)
	}
}

func TestClient_Models(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/models" {
			t.Errorf("Expected path /api/v1/models, got %v", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %v", r.Method)
		}
		
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test_token" {
			t.Errorf("Expected Bearer test_token, got %v", authHeader)
		}
		
		response := ModelsResponse{
			Data: []Model{
				{ID: "GigaChat", Object: "model", OwnedBy: "sberbank"},
			},
			Object: "list",
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	models, err := client.Models()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if len(models.Data) != 1 {
		t.Errorf("Expected 1 model, got %d", len(models.Data))
	}
	if models.Data[0].ID != "GigaChat" {
		t.Errorf("Expected model ID GigaChat, got %v", models.Data[0].ID)
	}
}

func TestClient_Models_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	_, err := client.Models()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	
	if _, ok := err.(*GigaChatError); !ok {
		t.Errorf("Expected GigaChatError, got %T", err)
	}
}

func TestClient_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/chat/completions" {
			t.Errorf("Expected path /api/v1/chat/completions, got %v", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %v", r.Method)
		}
		
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		
		if len(req.Messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(req.Messages))
		}
		if req.Messages[0].Content != "Hello" {
			t.Errorf("Expected message 'Hello', got %v", req.Messages[0].Content)
		}
		
		response := ChatResponse{
			Choices: []ChatChoice{
				{
					Index: 0,
					Message: Message{
						Role:    "assistant",
						Content: "Hi there!",
					},
					FinishReason: "stop",
				},
			},
			Model:   GigaChat,
			Created: time.Now().Unix(),
			Usage: Usage{
				PromptTokens:     5,
				CompletionTokens: 3,
				TotalTokens:      8,
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	messages := []Message{
		{Role: "user", Content: "Hello"},
	}
	
	resp, err := client.Chat(messages)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if len(resp.Choices) != 1 {
		t.Errorf("Expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "Hi there!" {
		t.Errorf("Expected response 'Hi there!', got %v", resp.Choices[0].Message.Content)
	}
}

func TestClient_Chat_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if req.Model != GigaChatPro {
			t.Errorf("Expected model %v, got %v", GigaChatPro, req.Model)
		}
		if req.Temperature == nil || *req.Temperature != 0.7 {
			t.Errorf("Expected temperature 0.7, got %v", req.Temperature)
		}
		if req.MaxTokens == nil || *req.MaxTokens != 100 {
			t.Errorf("Expected max tokens 100, got %v", req.MaxTokens)
		}
		
		response := ChatResponse{
			Choices: []ChatChoice{
				{Message: Message{Content: "Response"}},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	messages := []Message{{Role: "user", Content: "Test"}}
	
	_, err := client.Chat(
		messages,
		WithModel(GigaChatPro),
		WithTemperature(0.7),
		WithMaxTokens(100),
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_Chat_ValidationError(t *testing.T) {
	tm := createMockTokenManager("test_token")
	client := NewClient(tm)
	
	tests := []struct {
		name     string
		messages []Message
	}{
		{
			name:     "empty messages",
			messages: []Message{},
		},
		{
			name: "invalid role",
			messages: []Message{
				{Role: "invalid", Content: "test"},
			},
		},
		{
			name: "empty content",
			messages: []Message{
				{Role: "user", Content: ""},
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Chat(tt.messages)
			if err == nil {
				t.Fatal("Expected validation error, got nil")
			}
			if _, ok := err.(*ValidationError); !ok {
				t.Errorf("Expected ValidationError, got %T", err)
			}
		})
	}
}

func TestClient_ChatStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "text/event-stream" {
			t.Errorf("Expected Accept: text/event-stream")
		}
		
		w.Header().Set("Content-Type", "text/event-stream")
		
		events := []string{
			`data: {"choices":[{"delta":{"content":"Hello"}}]}`,
			`data: {"choices":[{"delta":{"content":" world"}}]}`,
			`data: [DONE]`,
		}
		
		for _, event := range events {
			w.Write([]byte(event + "\n\n"))
		}
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	messages := []Message{{Role: "user", Content: "Test"}}
	
	var receivedEvents []*ChatResponse
	var done bool
	
	err := client.ChatStream(messages, func(event *ChatResponse, isDone bool, err error) {
		if err != nil {
			t.Errorf("Stream callback error: %v", err)
		}
		if isDone {
			done = true
		} else if event != nil {
			receivedEvents = append(receivedEvents, event)
		}
	})
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !done {
		t.Error("Expected done signal")
	}
	if len(receivedEvents) != 2 {
		t.Errorf("Expected 2 events, got %d", len(receivedEvents))
	}
}

func TestClient_GenerateImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if req.FunctionCall != "auto" {
			t.Errorf("Expected function_call auto, got %v", req.FunctionCall)
		}
		
		response := ChatResponse{
			Choices: []ChatChoice{
				{
					Message: Message{
						Content: `<img src="file_123" />`,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	resp, err := client.GenerateImage("A beautiful sunset")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if len(resp.Choices) == 0 {
		t.Error("Expected at least one choice")
	}
}

func TestClient_GenerateImage_EmptyPrompt(t *testing.T) {
	tm := createMockTokenManager("test_token")
	client := NewClient(tm)
	
	_, err := client.GenerateImage("")
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestClient_DownloadImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/api/v1/files/") {
			t.Errorf("Expected files path, got %v", r.URL.Path)
		}
		if !strings.HasSuffix(r.URL.Path, "/content") {
			t.Errorf("Expected path to end with /content, got %v", r.URL.Path)
		}
		
		w.Header().Set("Content-Type", "application/jpg")
		w.Write([]byte("fake_image_data"))
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	imageData, err := client.DownloadImage("file_123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if imageData == "" {
		t.Error("Expected base64 encoded image data")
	}
}

func TestClient_DownloadImage_EmptyFileID(t *testing.T) {
	tm := createMockTokenManager("test_token")
	client := NewClient(tm)
	
	_, err := client.DownloadImage("")
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestClient_CreateImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/chat/completions") {
			response := ChatResponse{
				Choices: []ChatChoice{
					{
						Message: Message{
							Content: `<img src="file_123" />`,
						},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		} else if strings.Contains(r.URL.Path, "/files/") {
			w.Write([]byte("image_data"))
		}
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	result, err := client.CreateImage("A sunset")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if result.FileID != "file_123" {
		t.Errorf("Expected file ID file_123, got %v", result.FileID)
	}
	if result.Content == "" {
		t.Error("Expected image content")
	}
}

func TestValidateMessages(t *testing.T) {
	tests := []struct {
		name      string
		messages  []Message
		expectErr bool
	}{
		{
			name:      "empty messages",
			messages:  []Message{},
			expectErr: true,
		},
		{
			name: "valid messages",
			messages: []Message{
				{Role: "user", Content: "Hello"},
			},
			expectErr: false,
		},
		{
			name: "invalid role",
			messages: []Message{
				{Role: "invalid", Content: "test"},
			},
			expectErr: true,
		},
		{
			name: "empty content",
			messages: []Message{
				{Role: "user", Content: "   "},
			},
			expectErr: true,
		},
		{
			name: "multiple valid messages",
			messages: []Message{
				{Role: "system", Content: "You are helpful"},
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi"},
			},
			expectErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessages(tt.messages)
			if (err != nil) != tt.expectErr {
				t.Errorf("validateMessages() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestExtractImageID(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "valid img tag",
			content:  `<img src="file_123" />`,
			expected: "file_123",
		},
		{
			name:     "img tag with attributes",
			content:  `<img alt="test" src="file_456" class="image" />`,
			expected: "file_456",
		},
		{
			name:     "no img tag",
			content:  "Just text",
			expected: "",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractImageID(tt.content)
			if result != tt.expected {
				t.Errorf("extractImageID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAsk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ChatResponse{
			Choices: []ChatChoice{
				{Message: Message{Content: "Answer"}},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	tm := createMockTokenManager("test_token")
	client := NewClient(tm, WithBaseURI(server.URL))
	
	answer, err := Ask(client, "Question?")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if answer != "Answer" {
		t.Errorf("Expected 'Answer', got %v", answer)
	}
}

func TestConversation(t *testing.T) {
	messages := Conversation("You are helpful", "Hello")
	
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}
	if messages[0].Role != "system" {
		t.Errorf("Expected first role to be system, got %v", messages[0].Role)
	}
	if messages[1].Role != "user" {
		t.Errorf("Expected second role to be user, got %v", messages[1].Role)
	}
}

func TestConversation_EmptySystem(t *testing.T) {
	messages := Conversation("", "Hello")
	
	if len(messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(messages))
	}
	if messages[0].Role != "user" {
		t.Errorf("Expected role to be user, got %v", messages[0].Role)
	}
}

func TestContinueChat(t *testing.T) {
	existing := []Message{
		{Role: "user", Content: "First"},
	}
	
	continued := ContinueChat(existing, "Second")
	
	if len(continued) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(continued))
	}
	if continued[1].Content != "Second" {
		t.Errorf("Expected 'Second', got %v", continued[1].Content)
	}
}

func TestExtractContent(t *testing.T) {
	resp := &ChatResponse{
		Choices: []ChatChoice{
			{Message: Message{Content: "Test content"}},
		},
	}
	
	content := ExtractContent(resp)
	if content != "Test content" {
		t.Errorf("Expected 'Test content', got %v", content)
	}
}

func TestExtractContent_NoChoices(t *testing.T) {
	resp := &ChatResponse{
		Choices: []ChatChoice{},
	}
	
	content := ExtractContent(resp)
	if content != "" {
		t.Errorf("Expected empty string, got %v", content)
	}
}
