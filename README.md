# Go Clova

An unofficial Go library for the [Naver Clova Studio API](https://api.ncloud-docs.com/docs/en/ai-naver-clovastudio-summary).

This package has support for:
- Chat Completions

## Installation
```
go get github.com/starpia-forge/go-clova
```

## Usage

You can find the examples in the `examples` directory.

### Chat Completion:
```go
package main

import (
	"context"
	"fmt"
	"github.com/starpia-forge/go-clova"
	"os"
)

func main() {
	config := clova.DefaultConfig("Your Naver Clova Studio API Key")
	config.BaseURL = clova.NaverClovaAPIURLTestApp // for Test APP URL
	client := clova.NewClientWithConfig(config)

	response, err := client.CreateChatCompletion(context.Background(), clova.ModelHCXDASH001, clova.CompletionRequest{
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

	fmt.Println("Result Message Role :", response.Result.Message.Role)
	fmt.Println("Result Message Content :", response.Result.Message.Content)
}
```

### Chat Completion Streaming:
```go
package main

import (
	"context"
	"fmt"
	"github.com/starpia-forge/go-clova"
	"io"
	"os"
)

func main() {
	config := clova.DefaultConfig("Your Naver Clova Studio API Key")
	config.BaseURL = clova.NaverClovaAPIURLTestApp // for Test APP URL
	client := clova.NewClientWithConfig(config)

	stream, err := client.CreateChatCompletionStream(context.Background(), clova.ModelHCXDASH001, clova.CompletionRequest{
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
		fmt.Println("ID :", response.ID)
		fmt.Println("Event :", response.Event)
		fmt.Println("Data :", response.Data)
		fmt.Println()
	}
}
```

## Acknowledgment

This project was strongly inspired by [go-openai](https://github.com/sashabaranov/go-openai).

Therefore, while this project has different goals, it shares similar design principles with the original project.

Thank you to all the contributors who created the amazing `go-openai` library.

## License
[Apache License Version 2.0](https://github.com/StarpiaForge/go-clova/blob/master/LICENSE)