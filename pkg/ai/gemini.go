package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

type GeminiAIClient struct {
	client *genai.Client
	chat  *genai.Chat
	workspace *Workspace
}

func NewGeminiAIClient(apiKey string, workspace *Workspace) (*GeminiAIClient, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiAIClient{
		client: client,
		workspace: workspace,
	}, nil
}

func (g *GeminiAIClient) StartSession(ctx context.Context) (Response, error) {
	workspace := g.workspace
	session, err := workspace.GetSession("Session", "")

	if err != nil {
		return Response{}, fmt.Errorf("failed to start a session: %w", err)
	}

	instructions := fmt.Sprintf("%s\n%s", session.Role.Persona, workspace.Context.Settings.SystemPrompt)
	responseSchema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"think":   {Type: genai.TypeString},
			"summary": {Type: genai.TypeString},
			"content": {Type: genai.TypeString},
		},
		Required: []string{"think", "summary", "content"},
	}

	genConfig := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
		SystemInstruction: genai.NewContentFromText(instructions, genai.Role(session.Role.Name)),
	}

	g.chat, err = g.client.Chats.Create(ctx, "gemini-2.5-flash-preview-05-20", genConfig, nil)
	if err != nil {
		return Response{}, fmt.Errorf("failed to start a chat: %w", err)
	}
	var message strings.Builder
	if len(session.Chat) > 0 {
		message.WriteString("**Chat Context**: \n")
		for _, v := range session.Chat {
			message.WriteString(fmt.Sprintf("[user-message]: %s \n [agent-response]: %s", v.Message.Content, v.Response.Content))
		}
	} else {
		message.WriteString("Greetings")
	}

	return g.SendMessage(ctx, message.String(), nil, false)
}

func (g *GeminiAIClient) SendMessage(ctx context.Context, message string, history []Message, save bool) (Response, error) {
	if g.chat == nil {
		return Response{}, errors.New("chat session not started. Call StartSession first.")
	}

	resp, err := g.chat.SendMessage(ctx, genai.Part{
		Text: message,
	})

	if err != nil {
		return Response{}, fmt.Errorf("failed to get response from Gemini: %w", err)
	}

	if resp.Candidates == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return Response{}, errors.New("no response content received from Gemini model")
	}

	var responseText strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			responseText.WriteString(part.Text)
		}
	}

	rawAIResponse := responseText.String()
	respStruct, err := parseAIResponse(rawAIResponse)
	if err != nil {
		return Response{}, fmt.Errorf("failed to parse AI response into structured format: %w", err)
	}

	if _, err := g.workspace.GetActiveSession(); err == nil && save {
		g.workspace.AddInteraction(message, respStruct.Summary)
	}

	return respStruct, nil
}
