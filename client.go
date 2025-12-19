package gigachat

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	tokenManager *TokenManager
	baseURI      string
	httpClient   *http.Client
	defaultModel string
}

func NewClient(tokenManager *TokenManager, options ...ClientOption) *Client {
	c := &Client{
		tokenManager: tokenManager,
		baseURI:      "https://gigachat.devices.sberbank.ru",
		defaultModel: GigaChat,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

type ClientOption func(*Client)

func WithBaseURI(uri string) ClientOption {
	return func(c *Client) {
		c.baseURI = strings.TrimRight(uri, "/")
	}
}

func WithDefaultModel(model string) ClientOption {
	return func(c *Client) {
		c.defaultModel = model
	}
}

func WithClientInsecureSkipVerify(skip bool) ClientOption {
	return func(c *Client) {
		if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
			transport.TLSClientConfig.InsecureSkipVerify = skip
		}
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

func (c *Client) Models() (*ModelsResponse, error) {
	token, err := c.tokenManager.GetAccessToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.baseURI+"/api/v1/models", nil)
	if err != nil {
		return nil, &GigaChatError{Message: "failed to create request", Err: err}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &GigaChatError{Message: "request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &GigaChatError{
			Message: fmt.Sprintf("API request failed: %s", string(body)),
			Code:    resp.StatusCode,
		}
	}

	var modelsResp ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, &GigaChatError{Message: "failed to decode response", Err: err}
	}

	return &modelsResp, nil
}

func (c *Client) Chat(messages []Message, options ...ChatOption) (*ChatResponse, error) {
	if err := validateMessages(messages); err != nil {
		return nil, err
	}

	token, err := c.tokenManager.GetAccessToken()
	if err != nil {
		return nil, err
	}

	chatReq := ChatRequest{
		Model:    c.defaultModel,
		Messages: messages,
		Stream:   false,
	}

	for _, opt := range options {
		opt(&chatReq)
	}

	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return nil, &GigaChatError{Message: "failed to marshal request", Err: err}
	}

	req, err := http.NewRequest("POST", c.baseURI+"/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &GigaChatError{Message: "failed to create request", Err: err}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &GigaChatError{Message: "request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &GigaChatError{
			Message: fmt.Sprintf("API request failed: %s", string(body)),
			Code:    resp.StatusCode,
		}
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, &GigaChatError{Message: "failed to decode response", Err: err}
	}

	return &chatResp, nil
}

type StreamCallback func(event *ChatResponse, done bool, err error)

func (c *Client) ChatStream(messages []Message, callback StreamCallback, options ...ChatOption) error {
	if err := validateMessages(messages); err != nil {
		return err
	}

	token, err := c.tokenManager.GetAccessToken()
	if err != nil {
		return err
	}

	chatReq := ChatRequest{
		Model:    c.defaultModel,
		Messages: messages,
		Stream:   true,
	}

	for _, opt := range options {
		opt(&chatReq)
	}

	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		return &GigaChatError{Message: "failed to marshal request", Err: err}
	}

	req, err := http.NewRequest("POST", c.baseURI+"/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return &GigaChatError{Message: "failed to create request", Err: err}
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &GigaChatError{Message: "request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &GigaChatError{
			Message: fmt.Sprintf("API request failed: %s", string(body)),
			Code:    resp.StatusCode,
		}
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data:") {
			dataPart := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if dataPart == "[DONE]" {
				callback(nil, true, nil)
				return nil
			}

			var event ChatResponse
			if err := json.Unmarshal([]byte(dataPart), &event); err != nil {
				callback(nil, false, &GigaChatError{Message: "failed to decode event", Err: err})
				continue
			}

			callback(&event, false, nil)
		}
	}

	if err := scanner.Err(); err != nil {
		return &GigaChatError{Message: "stream reading error", Err: err}
	}

	return nil
}

func (c *Client) GenerateImage(prompt string, options ...ImageOption) (*ChatResponse, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, &ValidationError{Message: "image prompt cannot be empty"}
	}

	imgOpts := &imageOptions{}
	for _, opt := range options {
		opt(imgOpts)
	}

	messages := []Message{}
	if imgOpts.systemMessage != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: imgOpts.systemMessage,
		})
	}

	messages = append(messages, Message{
		Role:    "user",
		Content: prompt,
	})

	chatOpts := []ChatOption{WithFunctionCall("auto")}
	if imgOpts.model != "" {
		chatOpts = append(chatOpts, WithModel(imgOpts.model))
	}
	if imgOpts.temperature != nil {
		chatOpts = append(chatOpts, WithTemperature(*imgOpts.temperature))
	}

	return c.Chat(messages, chatOpts...)
}

func (c *Client) DownloadImage(fileID string) (string, error) {
	if strings.TrimSpace(fileID) == "" {
		return "", &ValidationError{Message: "file ID cannot be empty"}
	}

	token, err := c.tokenManager.GetAccessToken()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", c.baseURI+"/api/v1/files/"+fileID+"/content", nil)
	if err != nil {
		return "", &GigaChatError{Message: "failed to create request", Err: err}
	}

	req.Header.Set("Accept", "application/jpg")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", &GigaChatError{Message: "request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", &GigaChatError{
			Message: fmt.Sprintf("API request failed: %s", string(body)),
			Code:    resp.StatusCode,
		}
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", &GigaChatError{Message: "failed to read image data", Err: err}
	}

	return base64.StdEncoding.EncodeToString(imageData), nil
}

func (c *Client) CreateImage(prompt string, options ...ImageOption) (*ImageResult, error) {
	response, err := c.GenerateImage(prompt, options...)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) == 0 {
		return nil, &GigaChatError{Message: "no choices in response"}
	}

	content := response.Choices[0].Message.Content
	fileID := extractImageID(content)
	if fileID == "" {
		return nil, &GigaChatError{Message: "could not extract image ID from response"}
	}

	imageContent, err := c.DownloadImage(fileID)
	if err != nil {
		return nil, err
	}

	return &ImageResult{
		Content:  imageContent,
		FileID:   fileID,
		Response: *response,
	}, nil
}

type ChatOption func(*ChatRequest)

func WithModel(model string) ChatOption {
	return func(cr *ChatRequest) {
		cr.Model = model
	}
}

func WithTemperature(temp float64) ChatOption {
	return func(cr *ChatRequest) {
		cr.Temperature = &temp
	}
}

func WithTopP(topP float64) ChatOption {
	return func(cr *ChatRequest) {
		cr.TopP = &topP
	}
}

func WithMaxTokens(maxTokens int) ChatOption {
	return func(cr *ChatRequest) {
		cr.MaxTokens = &maxTokens
	}
}

func WithRepetitionPenalty(penalty float64) ChatOption {
	return func(cr *ChatRequest) {
		cr.RepetitionPenalty = &penalty
	}
}

func WithUpdateInterval(interval int) ChatOption {
	return func(cr *ChatRequest) {
		cr.UpdateInterval = &interval
	}
}

func WithFunctionCall(functionCall string) ChatOption {
	return func(cr *ChatRequest) {
		cr.FunctionCall = functionCall
	}
}

type imageOptions struct {
	systemMessage string
	model         string
	temperature   *float64
}

type ImageOption func(*imageOptions)

func WithSystemMessage(message string) ImageOption {
	return func(io *imageOptions) {
		io.systemMessage = message
	}
}

func WithImageModel(model string) ImageOption {
	return func(io *imageOptions) {
		io.model = model
	}
}

func WithImageTemperature(temp float64) ImageOption {
	return func(io *imageOptions) {
		io.temperature = &temp
	}
}

func validateMessages(messages []Message) error {
	if len(messages) == 0 {
		return &ValidationError{Message: "messages array cannot be empty"}
	}

	for i, msg := range messages {
		if msg.Role != "user" && msg.Role != "system" && msg.Role != "assistant" {
			return &ValidationError{
				Message: fmt.Sprintf("invalid role '%s' at index %d. Must be 'user', 'system', or 'assistant'", msg.Role, i),
			}
		}

		if strings.TrimSpace(msg.Content) == "" {
			return &ValidationError{
				Message: fmt.Sprintf("message content at index %d must be a non-empty string", i),
			}
		}
	}

	return nil
}

func extractImageID(content string) string {
	re := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func Ask(client *Client, question string, options ...ChatOption) (string, error) {
	messages := []Message{
		{Role: "user", Content: question},
	}

	response, err := client.Chat(messages, options...)
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", &GigaChatError{Message: "no choices in response"}
	}

	return response.Choices[0].Message.Content, nil
}

func Conversation(systemPrompt, userMessage string) []Message {
	messages := []Message{}
	if systemPrompt != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: userMessage,
	})
	return messages
}

func ContinueChat(messages []Message, userMessage string) []Message {
	return append(messages, Message{
		Role:    "user",
		Content: userMessage,
	})
}

func ExtractContent(response *ChatResponse) string {
	if len(response.Choices) == 0 {
		return ""
	}
	return response.Choices[0].Message.Content
}
