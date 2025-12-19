package gigachat

const (
	GigaChat2    = "GigaChat-2"
	GigaChat2Pro = "GigaChat-2-Pro"
	GigaChat2Max = "GigaChat-2-Max"
	GigaChat     = "GigaChat"
	GigaChatPro  = "GigaChat-Pro"
	GigaChatMax  = "GigaChat-Max"

	Embeddings      = "Embeddings"
	EmbeddingsGigaR = "EmbeddingsGigaR"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	Temperature       *float64  `json:"temperature,omitempty"`
	TopP              *float64  `json:"top_p,omitempty"`
	MaxTokens         *int      `json:"max_tokens,omitempty"`
	RepetitionPenalty *float64  `json:"repetition_penalty,omitempty"`
	UpdateInterval    *int      `json:"update_interval,omitempty"`
	Stream            bool      `json:"stream"`
	FunctionCall      string    `json:"function_call,omitempty"`
}

type ChatChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Delta        Message `json:"delta,omitempty"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Usage   Usage        `json:"usage"`
	Object  string       `json:"object"`
}

type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

type ModelsResponse struct {
	Data   []Model `json:"data"`
	Object string  `json:"object"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type ImageResult struct {
	Content  string       `json:"content"`
	FileID   string       `json:"file_id"`
	Response ChatResponse `json:"response"`
}

func GetGenerationModels() []string {
	return []string{
		GigaChat2,
		GigaChat2Pro,
		GigaChat2Max,
		GigaChat,
		GigaChatPro,
		GigaChatMax,
	}
}

func GetEmbeddingModels() []string {
	return []string{
		Embeddings,
		EmbeddingsGigaR,
	}
}

func IsValidGenerationModel(model string) bool {
	for _, m := range GetGenerationModels() {
		if m == model {
			return true
		}
	}
	return false
}

func IsValidEmbeddingModel(model string) bool {
	for _, m := range GetEmbeddingModels() {
		if m == model {
			return true
		}
	}
	return false
}
