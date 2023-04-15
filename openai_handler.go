package main

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

func initOpenAIClient() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	openaiClient = openai.NewClient(apiKey)
}

func handleAskCommand(question string) (string, error) {
	answer, err := callOpenAI(question)
	if err != nil {
		return "", err
	}

	if len(answer) > 2048 {
		answer = answer[:2045] + "..."
	}

	return answer, nil
}

func callOpenAI(question string) (string, error) {
	ctx := context.Background()

	content := "Question: " + question + "\nAnswer:"

	// Create a new message
	message := openai.ChatCompletionMessage{
		Role:    "user",
		Content: content,
	}

	// Configure the API request options
	options := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    []openai.ChatCompletionMessage{message},
		MaxTokens:   100,
		Temperature: 0.8,
		N:           1,
		TopP:        1,
	}

	// Call the OpenAI API
	completion, err := openaiClient.CreateChatCompletion(ctx, options)
	if err != nil {
		return "", err
	}

	// Get the answer from the response
	if len(completion.Choices) == 0 {
		return "You have literally rendered OpenAI speechless. Great work.", nil
	}

	return completion.Choices[0].Message.Content, nil
}
