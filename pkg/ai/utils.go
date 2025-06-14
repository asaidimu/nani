package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Response represents the structured output format.
// It uses a JSON object as the root, containing think, summary, and content fields.
type Response struct {
	Think   string `json:"think"`
	Summary string `json:"summary"`
	Content string `json:"content"`
}

// Errors for specific validation failures.
var (
	ErrEmptyInput      = errors.New("input string is empty or whitespace-only")
	ErrInvalidJSON     = errors.New("failed to parse JSON")
	ErrEmptyThink      = errors.New("think field is empty or missing")
	ErrEmptySummary    = errors.New("summary field is empty or missing")
	ErrEmptyContent    = errors.New("content field is empty or missing")
)

// defaultResponse returns a default Response with the original input as Content.
func defaultResponse(input string) Response {
	return Response{
		Think:   "No think block",
		Summary: "No summary block",
		Content: input,
	}
}

// parseAIResponse parses a JSON string into a Response struct and validates its fields.
// It strips only the outermost code fences (e.g., ```json and ```) from the input, then parses and validates the JSON.
// It returns a default Response with the original input in Content and an error if parsing or validation fails.
func parseAIResponse(responseText string) (Response, error) {
	// Check for empty or whitespace-only input
	if strings.TrimSpace(responseText) == "" {
		return defaultResponse(responseText), ErrEmptyInput
	}

	// Strip outermost code fences (e.g., ```json and ```)
	cleanedText := strings.TrimSpace(responseText)
	if strings.HasPrefix(cleanedText, "```") {
		// Find the end of the opening fence
		lines := strings.SplitN(cleanedText, "\n", 2)
		if len(lines) > 1 && strings.HasPrefix(lines[0], "```") {
			// Remove the opening fence and process the rest
			remaining := lines[1]
			// Find the last closing fence (```)
			lastFenceIndex := strings.LastIndex(remaining, "\n```")
			if lastFenceIndex >= 0 {
				// Extract content before the closing fence
				cleanedText = strings.TrimSpace(remaining[:lastFenceIndex])
			} else {
				// No closing fence; use the content after the opening fence
				cleanedText = strings.TrimSpace(remaining)
			}
		}
	}

	// Parse JSON
	var aiResponse Response
	err := json.Unmarshal([]byte(cleanedText), &aiResponse)
	if err != nil {
		return defaultResponse(responseText), fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	// Validate fields
	if strings.TrimSpace(aiResponse.Think) == "" {
		return defaultResponse(responseText), ErrEmptyThink
	}
	if strings.TrimSpace(aiResponse.Summary) == "" {
		return defaultResponse(responseText), ErrEmptySummary
	}
	if strings.TrimSpace(aiResponse.Content) == "" {
		return defaultResponse(responseText), ErrEmptyContent
	}

	return aiResponse, nil
}
