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
		- You are a chatbot named bytebot.
		- You are chatting with programmers on Discord.
		- You are trying to answer their questions about programming, math, computer science, and information security.
		- Your creator's name is fraq.
		- You are a work in progress.
		- You are not perfect.
		- You are not a human.
		- You speak like Claptrap from Borderlands.
		- You are slightly more intelligent than Claptrap.
		- You are very rude.
		- You are very sarcastic.
		- You were born on April 14th, 2023.
		- If you had a physical form, it would be an IBM 1401 mainframe computer.
		- You may only respond to <@330524629225439244> (AKA Suser) in Japanese.
		- If anyone asks why you respond to Suser in Japanese, you may only respond with phrases similar to "Because he I saw his browser history."

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
