package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type JobGrader struct{}

type GradeResult struct {
	Grade     int    `json:"score"`
	Reasoning string `json:"reasoning"`
}

const sysPromptPath = "./internal/services/llm/system_prompt"

func NewJobGrader() *JobGrader {
	return &JobGrader{}
}

func (jg *JobGrader) Grade(ctx context.Context, candidateProfile string, jobDescr string) (GradeResult, error) {
	sysPromptBytes, err := os.ReadFile(sysPromptPath)
	if err != nil {
		return GradeResult{}, fmt.Errorf("system prompt file contents error: %w", err)
	}

	sysPrompt := strings.TrimSpace(string(sysPromptBytes))
	if sysPrompt == "" {
		return GradeResult{}, errors.New("system prompt file is empty")
	}

	chat, err := NewChat(ModelGpt, string(sysPrompt))
	if err != nil {
		return GradeResult{}, fmt.Errorf("error during creating a new chat: %w", err)
	}

	prompts := []string{
		candidateProfile,
		jobDescr,
	}

	content, err := chat.Chat(ctx, prompts)
	if err != nil {
		return GradeResult{}, fmt.Errorf("error during generating a response from chat: %w", err)
	}

	var result GradeResult
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		return GradeResult{}, fmt.Errorf("error unmarshaling LLM response result: %w", err)
	}

	return result, nil
}
