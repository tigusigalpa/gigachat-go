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

	conversation := gigachat.Conversation(
		"Ты полезный помощник программиста на Go",
		"Как создать REST API в Go?",
	)

	response, err := client.Chat(conversation)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("Assistant:", gigachat.ExtractContent(response))

	conversation = append(conversation, gigachat.Message{
		Role:    "assistant",
		Content: response.Choices[0].Message.Content,
	})

	conversation = gigachat.ContinueChat(conversation, "А как добавить middleware для аутентификации?")

	response, err = client.Chat(conversation)
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	fmt.Println("\nAssistant:", gigachat.ExtractContent(response))
}
