package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	gigachat "github.com/tigusigalpa/gigachat-go"
)

func main() {
	authKey := os.Getenv("GIGACHAT_AUTH_KEY")
	if authKey == "" {
		clientID := os.Getenv("GIGACHAT_CLIENT_ID")
		clientSecret := os.Getenv("GIGACHAT_CLIENT_SECRET")
		if clientID == "" || clientSecret == "" {
			log.Fatal("Please set GIGACHAT_AUTH_KEY or GIGACHAT_CLIENT_ID and GIGACHAT_CLIENT_SECRET")
		}
		authKey = base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	}

	tokenManager := gigachat.NewTokenManager(
		authKey,
		gigachat.WithScope("GIGACHAT_API_PERS"),
	)

	client := gigachat.NewClient(
		tokenManager,
		gigachat.WithDefaultModel(gigachat.GigaChat2Pro),
	)

	messages := []gigachat.Message{
		{Role: "system", Content: "Ты эксперт по квантовой физике"},
		{Role: "user", Content: "Объясни принцип квантовой суперпозиции простыми словами"},
	}

	temp := 0.7
	maxTokens := 500
	topP := 0.9

	response, err := client.Chat(
		messages,
		gigachat.WithModel(gigachat.GigaChat2Pro),
		gigachat.WithTemperature(temp),
		gigachat.WithMaxTokens(maxTokens),
		gigachat.WithTopP(topP),
	)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Response:")
	fmt.Println(response.Choices[0].Message.Content)
	fmt.Printf("\nModel: %s\n", response.Model)
	fmt.Printf("Tokens: prompt=%d, completion=%d, total=%d\n",
		response.Usage.PromptTokens,
		response.Usage.CompletionTokens,
		response.Usage.TotalTokens,
	)
}
