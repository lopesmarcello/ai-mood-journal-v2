// Package ai deals with LLM integration and management
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lopesmarcello/ai-journal/dto"
	"github.com/sashabaranov/go-openai"
)

type AIClient struct {
	client       *openai.Client
	systemPrompt string
}

func NewAIClient(apiKey string, systemPrompt string) *AIClient {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"

	return &AIClient{
		client:       openai.NewClientWithConfig(config),
		systemPrompt: systemPrompt,
	}
}

func (c *AIClient) GenerateInsight(ctx context.Context, content string) (*dto.AIInsightResponse, error) {
	fullSystemPrompt := c.systemPrompt + "\n\nIMPORTANT: You must return ONLY a valid JSON object. " +
		"Format: {\"summary\": \"string\", \"themes\": [\"string\"], \"feelings\": [\"string\"], \"reflection\": \"string\"}"

	var lastErr error
	maxRetries := 3

	for i := range maxRetries {
		resp, err := c.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: "llama-3.1-8b-instant",
				Messages: []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleSystem, Content: fullSystemPrompt},
					{Role: openai.ChatMessageRoleUser, Content: content},
				},
				ResponseFormat: &openai.ChatCompletionResponseFormat{
					Type: openai.ChatCompletionResponseFormatTypeJSONObject,
				},
			})

		if err == nil {
			var insight dto.AIInsightResponse
			err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &insight)
			if err == nil {
				return &insight, nil
			}
			lastErr = fmt.Errorf("failed to parse AI JSON response: %w", err)
		} else {
			lastErr = err
		}

		waitTime := time.Duration(1<<i) * time.Second
		time.Sleep(waitTime)
	}
	return nil, fmt.Errorf("AI failed after %d retries: %v", maxRetries, lastErr)
}
