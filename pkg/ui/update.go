package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asaidimu/nani/pkg/ai"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	content = iota
	history
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd      tea.Cmd
		vpCmd      tea.Cmd
		spCmd      tea.Cmd
		previewVpCmd tea.Cmd
		cmds []tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.history, vpCmd = m.history.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)
	m.content, previewVpCmd = m.content.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout = m.calculateLayout(msg.Width-4, msg.Height-2)
		m.ready = true

		m.history.Width = m.layout.LeftWidth - HistoryStyle.GetHorizontalFrameSize()
		m.history.Height = m.layout.HistoryHeight - HistoryStyle.GetVerticalFrameSize()

		availableTextareaHeight := m.layout.InputHeight - PromptStyle.GetVerticalFrameSize() - 6
		if availableTextareaHeight < 1 {
			availableTextareaHeight = 1
		}
		if availableTextareaHeight > 10 {
			availableTextareaHeight = 10
		}

		m.textarea.SetWidth(m.layout.LeftWidth - PromptStyle.GetHorizontalFrameSize())
		m.textarea.SetHeight(availableTextareaHeight)

		m.content.Width = m.layout.RightWidth - PreviewStyle.GetHorizontalFrameSize()
		m.content.Height = m.layout.TotalHeight - PreviewStyle.GetVerticalFrameSize()

		m.updateHistoryContent()
		m.updatePreviewContent()

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "k":
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focused = (m.focused + 1) % 2
			return m, nil
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
	case tea.MouseMsg:
		var cmd tea.Cmd
		if m.focused == content {
			m.content, cmd = m.content.Update(msg)
		} else {
			m.history, cmd = m.history.Update(msg)
		}
		cmds = append(cmds, cmd)

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

	cmds = append(cmds, taCmd, vpCmd, spCmd, previewVpCmd)
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

	if len(m.messages) > 0 || content.Len() > 0 {
		content.WriteString("\n")
	}

	var spinnerLine string
	if m.loading {
		spinnerLine = AIMsgStyle.Render("AI: " + m.spinner.View() + " Thinking...")
	} else {
		spinnerLine = AIMsgStyle.Render("AI: ")
	}

	content.WriteString(spinnerLine)

	m.history.SetContent(content.String())
	m.history.GotoBottom()
}

func (m *Model) sendToAI(message string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		response, err := m.aiClient.SendMessage(ctx, message, m.messages, true)
		return AIResponseMsg{Content: response.Content, Think: response.Think, Summary: response.Summary, Err: err}
	}
}
