package ai

import (
	"context"
	"time"
)

// Message represents a chat message
type Message struct {
	Role    string
	Content string
	Time    time.Time
}

// AIClient interface for AI communication
type AIClient interface {
	StartSession(ctx context.Context) (Response, error)
	SendMessage(ctx context.Context, message string, history []Message, save bool) (Response, error)
}
