package ui

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if !m.ready {
		return "Initializing AI Chat Terminal..."
	}

	// Get the history content (which now includes the spinner area)
	historyText := m.viewport.View()

	// History section:
	historyContent := TitleStyle.Render("Chat History") + "\n\n" + historyText
	historySection := HistoryStyle.
		Width(m.layout.LeftWidth).
		Height(m.layout.HistoryHeight).
		Render(historyContent)

	// Input section:
	inputContent := TitleStyle.Render("Input") + "\n\n" +
		m.textarea.View() + "\n\n" +
		HelpStyle.Render("Enter: Send • Tab: Toggle Preview • Q/Ctrl+C: Quit")
	inputSection := PromptStyle.
		Width(m.layout.LeftWidth).
		Height(m.layout.InputHeight).
		Render(inputContent)

	// Preview section:
	previewContent := TitleStyle.Render("Preview") + "\n\n" + m.previewVP.View()
	previewSection := PreviewStyle.
		Width(m.layout.RightWidth).
		Height(m.layout.TotalHeight).
		Render(previewContent)

	// Combine left column (history + input) vertically.
	leftColumn := lipgloss.JoinVertical(lipgloss.Top, historySection, inputSection)

	// Combine everything horizontally.
	return lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, previewSection)
}

// updatePreviewContent prepares the styled content for the preview viewport
func (m *Model) updatePreviewContent() {
	if !m.ready {
		return
	}

	var rawPreviewContent string

	// Get the available width for content inside the preview box
	contentWidth := m.layout.RightWidth - PreviewStyle.GetHorizontalFrameSize()

	if len(m.messages) > 0 {
		var lastAIContentMsg string
		for i := len(m.messages) - 1; i >= 0; i-- {
			if m.messages[i].Role == "ai-content" {
				lastAIContentMsg = m.messages[i].Content
				break
			}
		}

		if lastAIContentMsg != "" {
			rendered, err := glamour.Render(lastAIContentMsg, "dark")
			if err != nil {
				rawPreviewContent += ErrorStyle.Render("Render Error: "+err.Error()) + "\n\n" +
					lipgloss.NewStyle().Width(contentWidth).Render(lastAIContentMsg)
			} else {
				rawPreviewContent += lipgloss.NewStyle().Width(contentWidth).Render(rendered)
			}
		}
	} else {
		welcomeText := "Welcome to AI Chat Terminal!\n\n" +
			"Features:\n" +
			"• Real-time markdown preview\n" +
			"• Responsive layout\n" +
			"• Beautiful terminal UI\n" +
			"• AI conversation history\n\n" +
			HelpStyle.Render("Start typing to see your message preview here.")
		rawPreviewContent = TitleStyle.Render("Preview Panel") + "\n\n" +
			lipgloss.NewStyle().Width(contentWidth).Render(welcomeText)
	}

	m.previewVP.SetContent(rawPreviewContent)
	m.previewVP.GotoTop()
}
