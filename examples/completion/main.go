package main

import (
	"context"
	"fmt"
	"github.com/starpia-forge/go-clova"
	"os"
)

func main() {
	config := clova.DefaultConfig(os.Getenv("NAVER_CLOVA_API_KEY"))
	config.BaseURL = clova.NaverClovaAPIURLTestApp
	client := clova.NewClientWithConfig(config)

	response, err := client.CreateChatCompletion(context.Background(), clova.ModelHCXDASH001, clova.CompletionRequest{
		Messages: []clova.CompletionMessage{
			{
				Role:    clova.CompletionMessageRoleUser,
				Content: "hello, world!",
			},
			{
				Role:    clova.CompletionMessageRoleAssistant,
				Content: "Hello there! Nice to see you. How can I assist you today?",
			},
			{
				Role:    clova.CompletionMessageRoleUser,
				Content: "What was the first sentence I said?",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("-- Created Completion Response")
	fmt.Println("Status Code :", response.Status.Code)
	fmt.Println("Status Message :", response.Status.Message)
	fmt.Println("Result Message Role :", response.Result.Message.Role)
	fmt.Println("Result Message Content :", response.Result.Message.Content)
}
