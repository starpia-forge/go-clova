package main

import (
	"context"
	"fmt"
	"github.com/StarpiaForge/go-clova"
	"io"
	"os"
)

func main() {
	config := clova.DefaultConfig(os.Getenv("NAVER_CLOVA_API_KEY"))
	config.BaseURL = clova.NaverClovaAPIURLTestApp
	client := clova.NewClientWithConfig(config)

	stream, err := client.CreateChatCompletionStream(context.Background(), clova.ModelHCXDASH001, clova.CompletionRequest{
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
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// Complete
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(response.Message.Content)
	}
}
