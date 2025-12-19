# GigaChat Go SDK Examples

This directory contains comprehensive examples demonstrating various features of the GigaChat Go SDK.

## 📋 Prerequisites

Before running the examples, set up your environment variables:

```bash
# Method 1: Using pre-encoded Authorization Key
export GIGACHAT_AUTH_KEY=your_base64_encoded_auth_key

# Method 2: Using Client ID and Client Secret
export GIGACHAT_CLIENT_ID=your_client_id
export GIGACHAT_CLIENT_SECRET=your_client_secret
```

## 🚀 Running Examples

### Basic Usage
Simple chat interaction and model listing.

```bash
cd basic
go run main.go
```

### Streaming
Real-time streaming responses from GigaChat.

```bash
cd streaming
go run main.go
```

### Image Generation
Generate and save images using GigaChat's text2image function.

```bash
cd image_generation
go run main.go
```

### Conversation
Multi-turn conversations with context management.

```bash
cd conversation
go run main.go
```

### Advanced
Advanced configuration with custom options and parameters.

```bash
cd advanced
go run main.go
```

## 📝 Example Descriptions

### 1. Basic (`basic/main.go`)
- Initialize client with token manager
- List available models
- Send simple chat message
- Display response and token usage

### 2. Streaming (`streaming/main.go`)
- Set up streaming chat
- Handle streaming events with callback
- Display real-time responses
- Handle completion and errors

### 3. Image Generation (`image_generation/main.go`)
- Generate image from text prompt
- Use system message for styling
- Download and save generated image
- Handle image-specific errors

### 4. Conversation (`conversation/main.go`)
- Create conversation with system prompt
- Maintain conversation context
- Continue multi-turn dialogue
- Extract content from responses

### 5. Advanced (`advanced/main.go`)
- Configure token manager with custom options
- Set up client with specific model
- Use advanced chat parameters
- Monitor token usage and model information

## 🔧 Customization

Each example can be customized by modifying:

- **Model**: Change `gigachat.GigaChat2Pro` to other available models
- **Temperature**: Adjust creativity (0.0 - 2.0)
- **Max Tokens**: Control response length
- **System Message**: Customize AI behavior

## 📖 Documentation

For more information, see:
- [Main README](../README.md)
- [Official GigaChat API Documentation](https://developers.sber.ru/docs/ru/gigachat/api/overview)

## 💡 Tips

1. **Start with Basic**: Run the basic example first to verify your setup
2. **Check Errors**: All examples include proper error handling
3. **Token Limits**: Be mindful of token usage in your account
4. **SSL Issues**: If you encounter SSL errors, you can use `WithInsecureSkipVerify(true)` option

## 🐛 Troubleshooting

### Authentication Errors
- Verify your `GIGACHAT_AUTH_KEY` or credentials are correct
- Check that your API key hasn't expired
- Ensure you have the correct scope (`GIGACHAT_API_PERS`)

### Connection Errors
- Check your internet connection
- Verify the API endpoint is accessible
- Try disabling SSL verification for testing (not recommended for production)

### Image Generation Issues
- Ensure your prompt contains "нарисуй" (draw) or similar command
- Check that image generation is available in your plan
- Verify file permissions for saving images

## 📞 Support

If you encounter issues:
1. Check the [main documentation](../README.md)
2. Review [official API docs](https://developers.sber.ru/docs/ru/gigachat/api/overview)
3. Open an issue on [GitHub](https://github.com/tigusigalpa/gigachat-go/issues)
