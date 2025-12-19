# 🚀 GigaChat Go SDK

![GigaChat Golang SDK](https://github.com/user-attachments/assets/5cc6724d-fd74-4908-ba5d-4e473bd6e885)

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigusigalpa/gigachat-go?style=flat-square)](https://goreportcard.com/report/github.com/tigusigalpa/gigachat-go)

Full-featured Go SDK for Sber GigaChat API integration. Simple, powerful, and idiomatic Go client for working with
GigaChat AI models, including streaming support and image generation.

**🌐 Language:** English | [Русский](README-ru.md)

## ✨ Features

- 🔌 **Simple Integration** with GigaChat API
- 🔐 **Automatic OAuth & Token Management** with thread-safe caching
- 🎯 **Support for All Models** (GigaChat-2, GigaChat-2-Pro, GigaChat-2-Max)
- 📝 **Conversation Support** with context management
- ⚡ **Streaming Support** for real-time responses
- 🎨 **Image Generation** with text2image function
- 🖼️ **Automatic Image Download** and base64 encoding
- 🎭 **Image Styling** via system prompts
- 🔧 **Functional Options Pattern** for clean API
- 🧪 **Comprehensive Examples** for all features
- 📚 **Full Documentation** with code samples
- 🛡️ **Type-Safe** with proper error handling

## 📦 Installation

```bash
go get github.com/tigusigalpa/gigachat-go
```

## ⚙️ Configuration

### 1. Getting Authorization Credentials

To work with GigaChat API, you need to obtain authorization credentials:

1. Register at [Sber AI Personal Account](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)
2. Create a project and get **Client ID** and **Client Secret**
3. Generate **Authorization Key** (Base64 of "Client ID:Client Secret")

> 💡 **Detailed Instructions
**: [Creating a Project and Getting Keys](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)

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

## 💡 Usage

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

## 🤖 Available Models

GigaChat supports several models for different tasks. The current list of models is available in
the [official documentation](https://developers.sber.ru/docs/ru/gigachat/models).

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

## 🔧 Generation Parameters

Available parameters for customizing generation:

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

## 🎨 Image Generation

GigaChat supports image generation using the built-in text2image function. To create images, use the verb "нарисуй" (
draw) in the prompt and the `function_call: auto` parameter.

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

> **Important**: For image generation, the prompt must contain the verb "нарисуй" (draw) or similar drawing commands.
> The API automatically determines the need to call the text2image function when the `function_call: auto` parameter is
> present.

## ⚠️ Error Handling

The SDK provides specialized error types for different error scenarios:

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

> 📖 **More about errors
**: [Official GigaChat API Documentation](https://developers.sber.ru/docs/ru/gigachat/api/errors-description)

## 📚 Examples

The repository includes comprehensive examples:

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

## 🔧 Configuration Options

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

## 🧪 Testing

```bash
# Run tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📖 Documentation

- [Official GigaChat API Documentation](https://developers.sber.ru/docs/ru/gigachat/api/overview)
- [GigaChat Models](https://developers.sber.ru/docs/ru/gigachat/models)
- [API Reference](https://developers.sber.ru/docs/ru/gigachat/api/reference)
- [Error Codes](https://developers.sber.ru/docs/ru/gigachat/api/errors-description)
- [Quick Start Guide](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Igor Sazonov**

- GitHub: [@tigusigalpa](https://github.com/tigusigalpa)
- Email: sovletig@gmail.com

## 🙏 Acknowledgments

- [Sber GigaChat](https://developers.sber.ru/docs/ru/gigachat/api/overview) for the amazing AI API
- The Go community for excellent tools and libraries

## 📊 Project Status

This project is actively maintained. If you encounter any issues or have suggestions,
please [open an issue](https://github.com/tigusigalpa/gigachat-go/issues).

---

Made with ❤️ by [Igor Sazonov](https://github.com/tigusigalpa)
