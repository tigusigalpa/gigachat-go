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

	messages := []gigachat.Message{
		{Role: "user", Content: "Напиши короткую историю о космосе"},
	}

	fmt.Println("Streaming response:")
	fmt.Println("---")

	err := client.ChatStream(messages, func(event *gigachat.ChatResponse, done bool, err error) {
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		if done {
			fmt.Println("\n---")
			fmt.Println("✅ Done!")
			return
		}

		if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
			fmt.Print(event.Choices[0].Delta.Content)
		}
	})

	if err != nil {
		log.Fatalf("Stream failed: %v", err)
	}
}
