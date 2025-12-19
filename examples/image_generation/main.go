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

	fmt.Println("Generating image...")

	result, err := client.CreateImage(
		"Нарисуй красивый закат над морем",
		gigachat.WithSystemMessage("Ты — художник-импрессионист"),
	)
	if err != nil {
		log.Fatalf("Image generation failed: %v", err)
	}

	fmt.Printf("Image generated successfully!\n")
	fmt.Printf("File ID: %s\n", result.FileID)

	imageData, err := base64.StdEncoding.DecodeString(result.Content)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	filename := "sunset.jpg"
	if err := os.WriteFile(filename, imageData, 0644); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	fmt.Printf("Image saved to: %s\n", filename)
}
