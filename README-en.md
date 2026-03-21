# GigaChat Go SDK

![GigaChat Golang SDK](https://i.postimg.cc/YCWrRjhK/528512211-5cc6724d-fd74-4908-ba5d-4e473bd6e885.jpg)

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigusigalpa/gigachat-go?style=flat-square)](https://goreportcard.com/report/github.com/tigusigalpa/gigachat-go)

Golang SDK for Sber GigaChat API. Supports streaming, image generation, and conversations.

**🌐 Language:** English | [Русский](README.md)

**📦 Package:** [pkg.go.dev/github.com/tigusigalpa/gigachat-go](https://pkg.go.dev/github.com/tigusigalpa/gigachat-go)

## About

Library for integrating GigaChat into Go applications. Suitable for chatbots, content generation, and image creation.

### Why this SDK?

- Works with Sber's GigaChat
- Streaming and large data support
- Automatic token management with thread-safe cache
- Error handling and type safety
- Functional options for configuration

### Use Cases

- **Chatbots** — Telegram, Discord, Slack. SDK maintains conversation context.
- **Text generation** — articles, product descriptions, newsletters.
- **Images** — generation from descriptions, artist styling.
- **Analytics** — document summarization, review processing.
- **Search** — semantic search via embeddings.
- **Education** — interactive systems, knowledge testing.

### Comparison

| Feature | GigaChat Go SDK | Alternatives |
|---------|-----------------|------------------|
| Russian Language | Native | Limited |
| Streaming | Yes | Partial |
| Image Generation | Built-in | Separate API |
| Token Management | Automatic | Manual |
| Type Safety | Yes | Partial |

## Features

- GigaChat API integration
- Automatic OAuth token management (thread-safe cache)
- All models: GigaChat-2, GigaChat-2-Pro, GigaChat-2-Max
- Conversations with context
- Streaming responses
- Image generation (text2image)
- Image download in base64
- Styling via system prompts
- Functional options
- Type safety and error handling

## Installation

```bash
go get github.com/tigusigalpa/gigachat-go
```

## Configuration

### 1. Getting Authorization Credentials

You need API keys:

1. Register at [Sber AI Personal Account](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)
2. Create a project and get **Client ID** and **Client Secret**
3. Generate **Authorization Key** (Base64 of "Client ID:Client Secret")

> Details: [Creating a Project and Getting Keys](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)

### 2. Environment Setup

Set up your environment variables:

```bash
# Method 1: Using ready Authorization Key
export GIGACHAT_AUTH_KEY=your_base64_encoded_auth_key

# Method 2: Using Client ID and Client Secret (will auto-generate auth_key)
export GIGACHAT_CLIENT_ID=your_client_id
export GIGACHAT_CLIENT_SECRET=your_client_secret

# Optional settings
export GIGACHAT_SCOPE=GIGACHAT_API_PERS
```

## Usage

### Basic Usage

```go
package main

import (
    "encoding/base64"
    "fmt"
    "log"
    "os"

    gigachat "github.com/tigusigalpa/gigachat-go"
)

func main() {
    // Create auth key from credentials
    authKey := base64.StdEncoding.EncodeToString(
        []byte(os.Getenv("GIGACHAT_CLIENT_ID") + ":" + os.Getenv("GIGACHAT_CLIENT_SECRET")),
    )

    // Or use pre-encoded key
    // authKey := os.Getenv("GIGACHAT_AUTH_KEY")

    // Create token manager
    tokenManager := gigachat.NewTokenManager(authKey)

    // Create client
    client := gigachat.NewClient(tokenManager)

    // Get available models
    models, err := client.Models()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Available models:")
    for _, model := range models.Data {
        fmt.Printf("- %s\n", model.ID)
    }

    // Send a message
    messages := []gigachat.Message{
        {Role: "user", Content: "Hello! Tell me a joke"},
    }

    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Choices[0].Message.Content)
}
```

### Working with Conversations

```go
// Create conversation with system prompt
conversation := gigachat.Conversation(
    "You are a helpful Go programming assistant",
    "How do I create a REST API in Go?",
)

response, err := client.Chat(conversation)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Assistant:", gigachat.ExtractContent(response))

// Continue the conversation
conversation = append(conversation, gigachat.Message{
    Role:    "assistant",
    Content: response.Choices[0].Message.Content,
})

conversation = gigachat.ContinueChat(conversation, "How do I add authentication middleware?")

response, err = client.Chat(conversation)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Assistant:", gigachat.ExtractContent(response))
```

### Streaming Requests

```go
messages := []gigachat.Message{
    {Role: "user", Content: "Write a long story about space"},
}

fmt.Println("Streaming response:")

err := client.ChatStream(messages, func(event *gigachat.ChatResponse, done bool, err error) {
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }

    if done {
        fmt.Println("\n✅ Done!")
        return
    }

    if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
        fmt.Print(event.Choices[0].Delta.Content)
    }
})

if err != nil {
    log.Fatal(err)
}
```

### Advanced Usage with Options

```go
// Configure token manager
tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithScope("GIGACHAT_API_PERS"),
    gigachat.WithInsecureSkipVerify(false), // Set to true to skip SSL verification
)

// Configure client
client := gigachat.NewClient(
    tokenManager,
    gigachat.WithDefaultModel(gigachat.GigaChat2Pro),
    gigachat.WithBaseURI("https://gigachat.devices.sberbank.ru"),
)

// Send message with options
messages := []gigachat.Message{
    {Role: "system", Content: "You are a quantum physics expert"},
    {Role: "user", Content: "Explain quantum superposition in simple terms"},
}

response, err := client.Chat(
    messages,
    gigachat.WithModel(gigachat.GigaChat2Pro),
    gigachat.WithTemperature(0.7),
    gigachat.WithMaxTokens(500),
    gigachat.WithTopP(0.9),
)

if err != nil {
    log.Fatal(err)
}

fmt.Println(response.Choices[0].Message.Content)
fmt.Printf("Tokens used: %d\n", response.Usage.TotalTokens)
```

## Available Models

Current list: [documentation](https://developers.sber.ru/docs/ru/gigachat/models).

### Text Generation Models

| Model              | Description                               | Use Cases                       |
|--------------------|-------------------------------------------|---------------------------------|
| **GigaChat-2**     | Base second-generation model              | General tasks, dialogues        |
| **GigaChat-2-Pro** | Advanced model with improved capabilities | Complex tasks, creative writing |
| **GigaChat-2-Max** | Maximum model for the most complex tasks  | Professional tasks, analysis    |

### Embedding Models

| Model               | Description                          | Use Cases                          |
|---------------------|--------------------------------------|------------------------------------|
| **Embeddings**      | Base model for vector representation | Semantic search, clustering        |
| **EmbeddingsGigaR** | Improved model for embeddings        | Accurate search, semantic analysis |

### Using Model Constants

```go
// Use constants for generation
response, err := client.Chat(
    messages,
    gigachat.WithModel(gigachat.GigaChat2Pro),
)

// Get list of available models
generationModels := gigachat.GetGenerationModels()
embeddingModels := gigachat.GetEmbeddingModels()

// Check model validity
if gigachat.IsValidGenerationModel(gigachat.GigaChat2) {
    // Model is valid for generation
}
```

## Generation Parameters

```go
response, err := client.Chat(
    messages,
    gigachat.WithModel(gigachat.GigaChat2Pro),      // Model to use
    gigachat.WithTemperature(0.7),                   // Creativity (0.0 - 2.0)
    gigachat.WithTopP(0.9),                          // Nucleus sampling (0.0 - 1.0)
    gigachat.WithMaxTokens(1000),                    // Maximum tokens
    gigachat.WithRepetitionPenalty(1.1),             // Repetition penalty (0.0 - 2.0)
    gigachat.WithUpdateInterval(0),                  // Update interval for streaming
)
```

## Image Generation

To create images, use the verb "нарисуй" (draw) in the prompt.

### Basic Usage

```go
// Simple image generation
result, err := client.CreateImage("Нарисуй красивый закат над морем")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Image generated! File ID: %s\n", result.FileID)

// Decode and save image
imageData, err := base64.StdEncoding.DecodeString(result.Content)
if err != nil {
    log.Fatal(err)
}

err = os.WriteFile("sunset.jpg", imageData, 0644)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Image saved to sunset.jpg")
```

### Generation with System Prompt (Styling)

```go
// Generate in the style of a specific artist
result, err := client.CreateImage(
    "Нарисуй розового кота",
    gigachat.WithSystemMessage("Ты — Василий Кандинский"),
)

// Generate in a specific style
result, err := client.CreateImage(
    "Нарисуй космический корабль",
    gigachat.WithSystemMessage("Ты — художник-концептуалист научной фантастики"),
    gigachat.WithImageTemperature(0.8),
)
```

### Separate Generation and Download

```go
// Generate image (returns response with image ID)
response, err := client.GenerateImage("Нарисуй футуристический город")
if err != nil {
    log.Fatal(err)
}

// Extract image ID from response
content := response.Choices[0].Message.Content
// Response contains: <img src="file-id" fuse="true"/>

// Download image separately
imageData, err := client.DownloadImage(fileID)
if err != nil {
    log.Fatal(err)
}

// Save file
decodedData, _ := base64.StdEncoding.DecodeString(imageData)
os.WriteFile("city.jpg", decodedData, 0644)
```

### Available Image Methods

| Method                              | Description                               | Returns           |
|-------------------------------------|-------------------------------------------|-------------------|
| `GenerateImage(prompt, options...)` | Generates image and returns API response  | `*ChatResponse`   |
| `DownloadImage(fileID)`             | Downloads image by ID                     | `string` (base64) |
| `CreateImage(prompt, options...)`   | Generates and downloads image in one call | `*ImageResult`    |

### Image Generation Options

```go
result, err := client.CreateImage(
    "Нарисуй дракона",
    gigachat.WithSystemMessage("Ты — художник фэнтези"),
    gigachat.WithImageModel(gigachat.GigaChat2Pro),
    gigachat.WithImageTemperature(0.8),
)
```

### Error Handling for Image Generation

```go
result, err := client.CreateImage("Нарисуй дракона")
if err != nil {
    switch e := err.(type) {
    case *gigachat.ValidationError:
        fmt.Printf("Validation error: %v\n", e)
    case *gigachat.GigaChatError:
        fmt.Printf("GigaChat API error (code %d): %v\n", e.Code, e)
    case *gigachat.AuthenticationError:
        fmt.Printf("Authentication error: %v\n", e)
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
    return
}

fmt.Printf("Image saved with ID: %s\n", result.FileID)
```

### Example Prompts for Generation

```go
// Good prompts (contain "нарисуй")
client.CreateImage("Нарисуй закат в горах")
client.CreateImage("Нарисуй портрет кота в стиле ренессанса")
client.CreateImage("Нарисуй абстрактную композицию")

// Styling via system_message
client.CreateImage(
    "Нарисуй цветы",
    gigachat.WithSystemMessage("Ты — Клод Моне, рисуешь в стиле импрессионизма"),
)

client.CreateImage(
    "Нарисуй робота",
    gigachat.WithSystemMessage("Ты — концепт-художник для научно-фантастических фильмов"),
)
```

> The prompt must contain "нарисуй" or a similar command. API calls text2image automatically with `function_call: auto`.

## Error Handling

SDK provides typed errors:

```go
response, err := client.Chat(messages)
if err != nil {
    switch e := err.(type) {
    case *gigachat.AuthenticationError:
        // Authentication errors (invalid keys, expired token)
        fmt.Printf("Authentication error: %v\n", e)
    case *gigachat.ValidationError:
        // Validation errors (invalid message format)
        fmt.Printf("Validation error: %v\n", e)
    case *gigachat.GigaChatError:
        // General GigaChat API errors
        fmt.Printf("GigaChat error (code %d): %v\n", e.Code, e)
    default:
        fmt.Printf("Unknown error: %v\n", err)
    }
    return
}
```

### GigaChat API Error Codes

#### 🔐 Authentication Errors (400-401)

| Code | HTTP | Description                                       | Solution                                                                       |
|------|------|---------------------------------------------------|--------------------------------------------------------------------------------|
| 1    | 400  | `scope data format invalid`                       | Check scope field format                                                       |
| 4    | 401  | `Can't decode 'Authorization' header`             | Check authorization key correctness                                            |
| 5    | 400  | `scope is empty`                                  | Specify scope: `GIGACHAT_API_PERS`, `GIGACHAT_API_B2B`, or `GIGACHAT_API_CORP` |
| 6    | 401  | `credentials doesn't match db data`               | Reissue authorization key in personal account                                  |
| 7    | 401  | `scope from db not fully includes consumed scope` | Specify correct API version in scope                                           |

#### 💳 Limit and Access Errors (402-403)

| HTTP | Description         | Cause                  | Solution                              |
|------|---------------------|------------------------|---------------------------------------|
| 402  | `Payment Required`  | Model tokens exhausted | Check token limit in personal account |
| 403  | `Permission denied` | No access to method    | Check tariff plan                     |

#### 📊 Data Size Errors (413)

| HTTP | Description         | Cause                    | Solution                   |
|------|---------------------|--------------------------|----------------------------|
| 413  | `Payload too large` | Input data size exceeded | Reduce prompt or file size |

#### ⚙️ Parameter Errors (422)

| HTTP | Description                                  | Cause                           | Solution                                 |
|------|----------------------------------------------|---------------------------------|------------------------------------------|
| 422  | `Requested model does not support functions` | Model doesn't support functions | Use different model or disable functions |
| 422  | `system message must be the first message`   | Incorrect message order         | System message must be first             |
| 422  | `Unprocessable Entity`                       | File exceeds context size       | Split or reduce file                     |

#### 🚦 Request Limit Errors (429)

| HTTP | Description         | Cause                             | Solution                             |
|------|---------------------|-----------------------------------|--------------------------------------|
| 429  | `Too Many Requests` | Concurrent request limit exceeded | Reduce request frequency, add delays |

#### 🔧 Server Errors (500)

| HTTP | Description             | Cause                  | Solution        |
|------|-------------------------|------------------------|-----------------|
| 500  | `Internal Server Error` | GigaChat service error | Contact support |

> Details: [GigaChat API Error Codes](https://developers.sber.ru/docs/ru/gigachat/api/errors-description)

## Examples

Repository includes examples:

- **[Basic Usage](examples/basic/main.go)** - Simple chat and model listing
- **[Streaming](examples/streaming/main.go)** - Real-time streaming responses
- **[Image Generation](examples/image_generation/main.go)** - Generate and save images
- **[Conversation](examples/conversation/main.go)** - Multi-turn conversations
- **[Advanced](examples/advanced/main.go)** - Advanced configuration and options

Run examples:

```bash
cd examples/basic
go run main.go

cd ../streaming
go run main.go

cd ../image_generation
go run main.go
```

## Configuration Options

### Token Manager Options

```go
tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithScope("GIGACHAT_API_PERS"),           // API scope
    gigachat.WithOAuthURI("https://..."),              // Custom OAuth URI
    gigachat.WithInsecureSkipVerify(true),             // Skip SSL verification
    gigachat.WithHTTPClient(customHTTPClient),         // Custom HTTP client
)
```

### Client Options

```go
client := gigachat.NewClient(
    tokenManager,
    gigachat.WithBaseURI("https://..."),               // Custom base URI
    gigachat.WithDefaultModel(gigachat.GigaChat2Pro),  // Default model
    gigachat.WithClientInsecureSkipVerify(true),       // Skip SSL verification
    gigachat.WithHTTPClient(customHTTPClient),         // Custom HTTP client
)
```

## Testing

The library includes a comprehensive set of unit tests for all major components.

### Test Structure

- **`errors_test.go`** - Tests for custom error types (`GigaChatError`, `AuthenticationError`, `ValidationError`)
- **`models_test.go`** - Tests for data models, constants, and JSON serialization
- **`token_manager_test.go`** - Tests for token management, caching, and concurrent access
- **`client_test.go`** - Tests for all API client methods using mock HTTP servers

### Running Tests

```bash
# Run all tests
go test -v

# Run specific test file
go test -v -run TestClient
go test -v -run TestTokenManager

# Run specific test
go test -v -run TestClient_Chat

# Run with code coverage
go test -cover

# Detailed coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run with race detector
go test -race -v
```

### Test Coverage

Tests cover:
- All public functions and methods
- Error handling and edge cases
- Input validation
- Concurrent token access
- JSON marshaling/unmarshaling
- Mock HTTP servers for isolation from real API

## Documentation

- [Official GigaChat API Documentation](https://developers.sber.ru/docs/ru/gigachat/api/overview)
- [GigaChat Models](https://developers.sber.ru/docs/ru/gigachat/models)
- [API Reference](https://developers.sber.ru/docs/ru/gigachat/api/reference)
- [Error Codes](https://developers.sber.ru/docs/ru/gigachat/api/errors-description)
- [Quick Start Guide](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)

## SDK Architecture

### Components

#### TokenManager
Manages OAuth authentication:
- Access token retrieval via OAuth 2.0
- Thread-safe token cache
- Auto-refresh on expiration
- Custom HTTP client support

#### Client
API interface:
- Chat, streaming, image generation
- Model and parameter management
- Typed errors
- Functional options

#### Models
Data structures:
- Model constants
- Parameter validation
- API response structures
- Message helpers

### Data Flow

```
Application → Client → TokenManager → OAuth Server
                ↓
           GigaChat API
                ↓
           Response → Client → Application
```

## Best Practices

### Managing Dialogue Context

```go
// Limit history size to save tokens
const maxHistoryMessages = 10

func trimHistory(messages []gigachat.Message) []gigachat.Message {
    if len(messages) <= maxHistoryMessages {
        return messages
    }
    // Keep system message (if any) and last N messages
    systemMsg := []gigachat.Message{}
    if messages[0].Role == "system" {
        systemMsg = messages[:1]
        messages = messages[1:]
    }
    return append(systemMsg, messages[len(messages)-maxHistoryMessages:]...)
}
```

### Handling Rate Limits

```go
import "time"

func chatWithRetry(client *gigachat.Client, messages []gigachat.Message, maxRetries int) (*gigachat.ChatResponse, error) {
    var response *gigachat.ChatResponse
    var err error
    
    for i := 0; i < maxRetries; i++ {
        response, err = client.Chat(messages)
        if err == nil {
            return response, nil
        }
        
        if gigachatErr, ok := err.(*gigachat.GigaChatError); ok {
            if gigachatErr.Code == 429 {
                // Exponential backoff on rate limit
                backoff := time.Duration(1<<uint(i)) * time.Second
                time.Sleep(backoff)
                continue
            }
        }
        return nil, err
    }
    return nil, err
}
```

### Optimizing Token Usage

```go
// Use max_tokens parameter to control costs
response, err := client.Chat(
    messages,
    gigachat.WithMaxTokens(500),  // Limit response length
    gigachat.WithTemperature(0.3), // Less creativity = fewer tokens
)

// Monitor token usage
if err == nil {
    log.Printf("Tokens used: %d (prompt: %d, completion: %d)",
        response.Usage.TotalTokens,
        response.Usage.PromptTokens,
        response.Usage.CompletionTokens,
    )
}
```

### Secure Credential Storage

```go
// ❌ Bad: hardcoded in code
authKey := "hardcoded_key_here"

// ✅ Good: using environment variables
authKey := os.Getenv("GIGACHAT_AUTH_KEY")
if authKey == "" {
    log.Fatal("GIGACHAT_AUTH_KEY not set")
}

// ✅ Even better: using secret managers
// E.g., HashiCorp Vault, AWS Secrets Manager, etc.
```

### Graceful Shutdown for Streaming

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    cancel()
}()

err := client.ChatStream(messages, func(event *gigachat.ChatResponse, done bool, err error) {
    select {
    case <-ctx.Done():
        return
    default:
        // Process event
    }
})
```

## Troubleshooting

### Issue: Authentication Error

**Symptoms:**
```
Authentication error: Can't decode 'Authorization' header
```

**Solution:**
1. Check Authorization Key format correctness (must be Base64)
2. Ensure you're using correct Client ID and Client Secret
3. Verify keys don't contain extra spaces or characters

```go
// Correct auth key creation
authKey := base64.StdEncoding.EncodeToString(
    []byte(strings.TrimSpace(clientID) + ":" + strings.TrimSpace(clientSecret)),
)
```

### Issue: SSL/TLS Errors

**Symptoms:**
```
x509: certificate signed by unknown authority
```

**Solution:**
```go
// For development (NOT for production!)
tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithInsecureSkipVerify(true),
)

// For production: install correct certificates
```

### Issue: Token Limit Exceeded

**Symptoms:**
```
HTTP 402: Payment Required
```

**Solution:**
1. Check token balance in personal account
2. Optimize prompts to reduce consumption
3. Use `max_tokens` parameter to limit responses
4. Consider upgrading to a higher pricing plan

### Issue: Slow Responses

**Solution:**
1. Use streaming for large responses
2. Choose appropriate model (GigaChat-2 is faster than GigaChat-2-Max)
3. Reduce prompt size and dialogue history
4. Check network connection

```go
// Use streaming for long responses
err := client.ChatStream(messages, func(event *gigachat.ChatResponse, done bool, err error) {
    // Receive response parts as they're generated
})
```

## FAQ

### Which model is best for my task?

- **GigaChat-2**: Suitable for simple dialogues, FAQ bots, basic text generation
- **GigaChat-2-Pro**: Optimal for complex tasks, creative writing, text analysis
- **GigaChat-2-Max**: Use for professional tasks requiring maximum accuracy

### Can I use the SDK in production?

Yes:
- Error handling
- Thread safety
- Automatic token management
- Tests

### How to limit token consumption?

```go
response, err := client.Chat(
    messages,
    gigachat.WithMaxTokens(500),      // Maximum tokens in response
    gigachat.WithTemperature(0.3),    // Less variability
)
```

### Is proxy support available?

Yes, use a custom HTTP client:

```go
proxyURL, _ := url.Parse("http://proxy:8080")
httpClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}

tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithHTTPClient(httpClient),
)
```

### How to handle large files?

Split files into chunks and process sequentially:

```go
const chunkSize = 4000 // characters

func processLargeText(client *gigachat.Client, text string) ([]string, error) {
    var results []string
    
    for i := 0; i < len(text); i += chunkSize {
        end := i + chunkSize
        if end > len(text) {
            end = len(text)
        }
        
        chunk := text[i:end]
        messages := []gigachat.Message{
            {Role: "user", Content: "Summarize text: " + chunk},
        }
        
        response, err := client.Chat(messages)
        if err != nil {
            return nil, err
        }
        
        results = append(results, gigachat.ExtractContent(response))
    }
    
    return results, nil
}
```

### Can I use multiple models simultaneously?

Yes, create multiple clients or use options:

```go
// Option 1: Different clients
client1 := gigachat.NewClient(tm, gigachat.WithDefaultModel(gigachat.GigaChat2))
client2 := gigachat.NewClient(tm, gigachat.WithDefaultModel(gigachat.GigaChat2Pro))

// Option 2: One client with options
response1, _ := client.Chat(messages, gigachat.WithModel(gigachat.GigaChat2))
response2, _ := client.Chat(messages, gigachat.WithModel(gigachat.GigaChat2Pro))
```

## Performance

### Parallel Request Processing

```go
func processParallel(client *gigachat.Client, prompts []string) []string {
    results := make([]string, len(prompts))
    var wg sync.WaitGroup
    
    for i, prompt := range prompts {
        wg.Add(1)
        go func(idx int, p string) {
            defer wg.Done()
            
            messages := []gigachat.Message{{Role: "user", Content: p}}
            response, err := client.Chat(messages)
            if err != nil {
                log.Printf("Error for prompt %d: %v", idx, err)
                return
            }
            results[idx] = gigachat.ExtractContent(response)
        }(i, prompt)
    }
    
    wg.Wait()
    return results
}
```

### Connection Pool

```go
// Configure HTTP client with connection pool
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
    Timeout: 30 * time.Second,
}

tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithHTTPClient(httpClient),
)
```

## Security

1. **Never commit credentials to Git**
   ```bash
   # Add to .gitignore
   .env
   *.key
   credentials.json
   ```

2. **Use environment variables**
   ```go
   authKey := os.Getenv("GIGACHAT_AUTH_KEY")
   ```

3. **Key rotation**
   Regularly update Client ID and Client Secret in personal account

4. **User input validation**
   ```go
   func sanitizeInput(input string) string {
       // Remove potentially dangerous characters
       input = strings.TrimSpace(input)
       if len(input) > 10000 {
           input = input[:10000]
       }
       return input
   }
   ```

5. **Application-side rate limiting**
   ```go
   import "golang.org/x/time/rate"
   
   limiter := rate.NewLimiter(rate.Every(time.Second), 5) // 5 requests per second
   
   func rateLimitedChat(client *gigachat.Client, messages []gigachat.Message) (*gigachat.ChatResponse, error) {
       if err := limiter.Wait(context.Background()); err != nil {
           return nil, err
       }
       return client.Chat(messages)
   }
   ```

## Contributing

Pull requests are welcome.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Guidelines

- Use `gofmt`
- Add tests
- Update documentation when changing API
- Clear commit messages
- All tests must pass

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Igor Sazonov**

- GitHub: [@tigusigalpa](https://github.com/tigusigalpa)
- Email: sovletig@gmail.com

## Acknowledgments

- [Sber GigaChat](https://developers.sber.ru/docs/ru/gigachat/api/overview)
- Go community

## Status

Project is maintained. Issues and suggestions — [open an issue](https://github.com/tigusigalpa/gigachat-go/issues).
