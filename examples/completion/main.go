package main

import (
	"context"
	"fmt"
	"github.com/StarpiaForge/go-clova"
	"os"
)

func main() {
	config := clova.DefaultConfig(os.Getenv("NAVER_CLOVA_API_KEY"))
	config.BaseURL = clova.NaverClovaAPIURLTestApp
	client := clova.NewClientWithConfig(config)

	response, err := client.CreateCompletion(context.Background(), clova.ModelHCXDASH001, clova.CompletionRequest{
		Messages: []clova.CompletionMessage{
			{
				Role:    clova.CompletionMessageRoleUser,
				Content: "hello, world!",
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
