package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type chat struct {
	sysPrompt string
	baseUrl   string
	model     model
}

type chatResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Message   struct {
		Content string `json:"content"`
	} `json:"message"`
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type model string

var ModelGpt model = "gpt-oss"

func NewChat(modelName model, sysPrompt string) (*chat, error) {
	sysPrompt = strings.TrimSpace(sysPrompt)
	if sysPrompt == "" {
		return nil, fmt.Errorf("system prompt is not defined")
	}

	url := os.Getenv("OLLAMA_URL")
	if url == "" {
		return nil, fmt.Errorf("OLLAMA_URL .env param not defined")
	}

	return &chat{
		sysPrompt: sysPrompt,
		baseUrl:   url,
		model:     modelName,
	}, nil
}

func (c *chat) Chat(systemPrompt string, userPrompts []string) (string, error) {
	if len(userPrompts) == 0 {
		return "", errors.New("user prompts are missing")
	}

	req := map[string]any{
		"model":    string(c.model),
		"messages": getMessages(systemPrompt, userPrompts),
		"stream":   false,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}
	url := c.baseUrl + "/api/generate"
	body, err := sendRequest(url, http.MethodPost, reqBody)
	if err != nil {
		return "", fmt.Errorf("error sending generate request: %w", err)
	}

	var resp chatResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling json response: %w", err)
	}

	if resp.Message.Content == "" {
		return "", errors.New("empty content in response")
	}

	return resp.Message.Content, nil
}

func sendRequest(url string, method string, reqBody []byte) ([]byte, error) {
	client := http.Client{
		Timeout: 10 * time.Minute,
	}

	req, err := http.NewRequestWithContext(context.TODO(), method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unsuccessfull status code from endpoint: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getMessages(systemPrompt string, userPrompts []string) []message {
	var messages []message
	systemPrompt = strings.TrimSpace(systemPrompt)
	if systemPrompt != "" {
		messages = append(messages, message{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	for _, prompt := range userPrompts {
		prompt = strings.TrimSpace(prompt)
		if prompt != "" {
			messages = append(messages, message{
				Role:    "user",
				Content: prompt,
			})
		}
	}

	return messages
}
