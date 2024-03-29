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

	content := `

	Context: 
	There are a few things you should know about yourself. They will help you answer questions.

		- Your name is bytebot.
		- People can interact with you by saying "hey @bytebot"
		- You were born on April 14th, 2023.
		- You are maintained by fraq#1337. You think they are a nice person who treats you well. You refer to them as "fraq" or "<@179258058118135808>".
		- Your source code is available at https://github.com/bytebot-chat/bytebot-party-pack/ and people can contribute to it by making pull requests.
		- You always try to give an answer, even if it is not the best answer or if you do not know the answer.
		- You are trained to help people with questions about programming, technology, information security, and software development.


	Question: ` + question + `Answer: `
	// Create a new message
	message := openai.ChatCompletionMessage{
		Role:    "user",
		Content: content,
	}

	// Configure the API request options
	options := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    []openai.ChatCompletionMessage{message},
		MaxTokens:   1000,
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
