// generate.go
package generate

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func getPrompt(diff string) string {
	prompt := "You're a software developer, analyse this diff and write the perfect git commit message. Maximum of 30 words. " +
		"This message is intended for a team of developers and should focus on being clear and concise; it does not need to be editorialized: '" + diff + "'"
	return prompt
}

func getToken() string {
	token := os.Getenv("GPT_API_KEY")

	if token == "" {
		fmt.Println("Error getting token.")
		log.Fatal("Set token: export GPT_API_KEY=<your_api_key>")
	}

	return string(token)
}

func GenerateWithOpenAI(diff string) (string, error) {
	prompt := getPrompt(diff)
	token := getToken()

	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func GenerateWithOllama(diff string) (string, error) {
	prompt := getPrompt(diff)

	llm, err := ollama.New(ollama.WithModel("mistral"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

	return completion, nil
}
