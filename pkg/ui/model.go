package ui

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/asaidimu/nani/pkg/ai"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Layout struct {
	LeftWidth     int
	RightWidth    int
	HistoryHeight int
	InputHeight   int
	TotalHeight   int
}

type Model struct {
	messages    []ai.Message
	textarea    textarea.Model
	history    viewport.Model
	content   viewport.Model
	spinner     spinner.Model
	loading     bool
	ready       bool
	aiClient    ai.AIClient
	layout      Layout
	previewMode bool
	focused     int
}

type AIResponseMsg struct {
	Content string
	Think string
	Summary string
	Err     error
}

type ErrMsg error

func New(aiClient ai.AIClient) *Model {
	ta := textarea.New()
	ta.Placeholder = "Type your message here... (Press Enter to send, Tab to toggle preview)"
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 2000000
	ta.SetHeight(3)
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	vp := viewport.New(50, 20)
	vp.KeyMap.Down.SetKeys("down", "pgdown")
	vp.KeyMap.Up.SetKeys("up", "pgup")

	previewVp := viewport.New(50, 20)
	previewVp.KeyMap.Up.SetKeys("up", "pgup")
	previewVp.KeyMap.Down.SetKeys("down", "pgdown")

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ctx := context.Background();
	response, err := aiClient.StartSession(ctx)

	if err != nil {
		fmt.Printf("Error initializing Gemini client: %v\n", err)
		os.Exit(1)
	}

	result := &Model{
		messages:    []ai.Message{},
		textarea:    ta,
		history:    vp,
		content:   previewVp,
		spinner:     s,
		aiClient:    aiClient,
		ready:       false,
		previewMode: false,
	}
	result.messages = append(result.messages, ai.Message{
		Role: "ai-content",
		Content: response.Content,
		Time: time.Now(),
	})
	return result
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

func (m *Model) calculateLayout(width, height int) Layout {
	minOverallWidth := 80
	minOverallHeight := 15

	if width < minOverallWidth {
		width = minOverallWidth
	}
	if height < minOverallHeight {
		height = minOverallHeight
	}

	leftWidth := int(float64(width) * 0.4)
	minColumnContentWidth := 20
	if leftWidth < minColumnContentWidth+HistoryStyle.GetHorizontalFrameSize() {
		leftWidth = minColumnContentWidth + HistoryStyle.GetHorizontalFrameSize()
	}
	rightWidth := width - leftWidth

	if rightWidth < minColumnContentWidth+PreviewStyle.GetHorizontalFrameSize() {
		rightWidth = minColumnContentWidth + PreviewStyle.GetHorizontalFrameSize()
		leftWidth = width - rightWidth
		if leftWidth < minColumnContentWidth+HistoryStyle.GetHorizontalFrameSize() {
			leftWidth = minColumnContentWidth + HistoryStyle.GetHorizontalFrameSize()
		}
	}

	minInputHeight := 8
	maxInputHeight := 15
	minHistoryHeight := 6

	proposedInputHeight := int(float64(height) * 0.25)

	inputHeight := proposedInputHeight
	if inputHeight < minInputHeight {
		inputHeight = minInputHeight
	}
	if inputHeight > maxInputHeight {
		inputHeight = maxInputHeight
	}

	historyHeight := height - inputHeight

	if historyHeight < minHistoryHeight {
		historyHeight = minHistoryHeight
		inputHeight = height - historyHeight
		if inputHeight < minInputHeight {
			inputHeight = minInputHeight
		}
	}

	return Layout{
		LeftWidth:     leftWidth,
		RightWidth:    rightWidth,
		HistoryHeight: historyHeight - 2,
		InputHeight:   inputHeight,
		TotalHeight:   height,
	}
}
