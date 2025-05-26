package ai

import (
	"context"
	"encoding/xml"
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
	StartSession(ctx context.Context) (string, error)
	SendMessage(ctx context.Context, message string, history []Message) (Response, error)
}

// AIResponseStructure represents the structured output format.
// It uses <response> as the root element, containing <think>, <summary>, and <content>.
type Response struct {
	XMLName xml.Name `xml:"response"` // The root element is now <response>

	// Think contains the detailed reasoning and thought process,
	// typically formatted in Markdown.
	// will be shown in chat history
	Think string `xml:"think"`

	// Summary provides a concise yet comprehensive overview of the request,
	// context, and what has been accomplished, in plain text.
	// Will also be shown in chat history
	Summary string `xml:"summary"`

	// Content holds the complete, detailed answer or solution,
	// formatted in Markdown, which can include code blocks.
	// Renamed from 'Response' to 'Content' to match the new structure.
	// Will also be shown in preview
	Content string `xml:"content"`
}
