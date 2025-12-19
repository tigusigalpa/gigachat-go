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

	tokenManager := gigachat.NewTokenManager(authKey)

	client := gigachat.NewClient(tokenManager)

	models, err := client.Models()
	if err != nil {
		log.Fatalf("Failed to get models: %v", err)
	}

	fmt.Println("Available models:")
	for _, model := range models.Data {
		fmt.Printf("- %s (owned by: %s)\n", model.ID, model.OwnedBy)
	}

	messages := []gigachat.Message{
		{Role: "user", Content: "Привет! Расскажи короткий анекдот"},
	}

	response, err := client.Chat(messages)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("\nResponse:")
	fmt.Println(response.Choices[0].Message.Content)
	fmt.Printf("\nTokens used: %d\n", response.Usage.TotalTokens)
}
