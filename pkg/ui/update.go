package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asaidimu/nani/pkg/ai"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd      tea.Cmd
		vpCmd      tea.Cmd
		spCmd      tea.Cmd
		previewVpCmd tea.Cmd
	)

	// Update components
	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)
	m.previewVP, previewVpCmd = m.previewVP.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout = m.calculateLayout(msg.Width-4, msg.Height-2)
		m.ready = true

		// Update component dimensions for their *content areas*
		// History viewport content area
		m.viewport.Width = m.layout.LeftWidth - HistoryStyle.GetHorizontalFrameSize()
		m.viewport.Height = m.layout.HistoryHeight - HistoryStyle.GetVerticalFrameSize()

		// Calculate available height for textarea within the input section
		// Input section contains: title (1 line) + textarea + help (1 line) + newlines (4 lines) + borders/padding
		availableTextareaHeight := m.layout.InputHeight - PromptStyle.GetVerticalFrameSize() - 6 // 6 = title + help + newlines
		if availableTextareaHeight < 1 {
			availableTextareaHeight = 1 // Minimum 1 line for textarea
		}
		if availableTextareaHeight > 10 {
			availableTextareaHeight = 10 // Maximum to keep reasonable
		}

		// Update textarea dimensions
		m.textarea.SetWidth(m.layout.LeftWidth - PromptStyle.GetHorizontalFrameSize())
		m.textarea.SetHeight(availableTextareaHeight)

		// Preview viewport content area
		m.previewVP.Width = m.layout.RightWidth - PreviewStyle.GetHorizontalFrameSize()
		m.previewVP.Height = m.layout.TotalHeight - PreviewStyle.GetVerticalFrameSize()

		// Re-render content with new sizes
		m.updateHistoryContent()
		m.updatePreviewContent()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.previewMode = !m.previewMode
			m.updatePreviewContent()
		case "enter":
			if !m.loading && m.textarea.Value() != "" {
				userMsg := strings.TrimSpace(m.textarea.Value())
				m.messages = append(m.messages, ai.Message{
					Role:    "user",
					Content: userMsg,
					Time:    time.Now(),
				})

				m.textarea.Reset()
				m.loading = true
				m.updateHistoryContent()
				m.updatePreviewContent()

				return m, tea.Batch(
					m.sendToAI(userMsg),
					m.spinner.Tick,
				)
			}
		}

	case AIResponseMsg:
		m.loading = false
		if msg.Err != nil {
			m.messages = append(m.messages, ai.Message{
				Role:    "ai-content",
				Content: msg.Content,
				Time:    time.Now(),
			})
		} else {
			m.messages = append(m.messages, ai.Message{
				Role:    "assistant",
				Content: fmt.Sprintf("Summary: %s\n\nThought Process: %s", msg.Summary, msg.Think), // Combine for history
				Time:    time.Now(),
			})

			m.messages = append(m.messages, ai.Message{
				Role:    "ai-content",
				Content: msg.Content,
				Time:    time.Now(),
			})
		}
		m.updateHistoryContent()
		m.updatePreviewContent()

	case ErrMsg:
		m.loading = false
		return m, nil
	}

	return m, tea.Batch(taCmd, vpCmd, spCmd, previewVpCmd)
}

func (m *Model) updateHistoryContent() {
	if !m.ready {
		return
	}

	var content strings.Builder

	// Get the available width for text content inside the history box
	contentWidth := m.layout.LeftWidth - HistoryStyle.GetHorizontalFrameSize()

	for i, msg := range m.messages {
		if i > 0 {
			content.WriteString("\n") // Add a newline between messages
		}

		var styledLine string
		if msg.Role == "user" {
			styledLine = UserMsgStyle.Width(contentWidth).Render("You: " + msg.Content)
		} else if msg.Role == "assistant" { // This will now show summary and think
			styledLine = AIMsgStyle.Width(contentWidth).Render("AI: " + msg.Content)
		} else if msg.Role == "ai-content" { // This message is for preview only, skip for history
			continue
		}
		content.WriteString(styledLine)
	}

	// Always add the spinner area to prevent UI shifts
	if len(m.messages) > 0 || content.Len() > 0 {
		content.WriteString("\n") // Add a newline before the spinner area
	}

	// Always reserve space for the spinner line, show spinner when loading or empty space when not
	var spinnerLine string
	if m.loading {
		spinnerLine = AIMsgStyle.Render("AI: " + m.spinner.View() + " Thinking...")
	} else {
		spinnerLine = AIMsgStyle.Render("AI: ") // Just the prefix to maintain consistent spacing
	}
	content.WriteString(spinnerLine)

	m.viewport.SetContent(content.String())
	m.viewport.GotoBottom()
}

func (m *Model) sendToAI(message string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		response, err := m.aiClient.SendMessage(ctx, message, m.messages)
		return AIResponseMsg{Content: response.Content, Think: response.Think, Summary: response.Summary, Err: err}
	}
}
