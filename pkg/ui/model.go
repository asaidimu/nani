package ui

import (
	"github.com/asaidimu/nani/pkg/ai"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Layout holds calculated dimensions
type Layout struct {
	LeftWidth     int
	RightWidth    int
	HistoryHeight int // Total height for the history lipgloss.Style box
	InputHeight   int // Total height for the prompt lipgloss.Style box
	TotalHeight   int // Total height of the terminal window
}

// Model holds the application state
type Model struct {
	messages    []ai.Message
	textarea    textarea.Model
	viewport    viewport.Model   // For chat history content
	previewVP   viewport.Model   // For preview panel content
	spinner     spinner.Model
	loading     bool
	ready       bool
	aiClient    ai.AIClient
	layout      Layout
	previewMode bool
}

// Messages for Bubble Tea
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
	ta.SetHeight(3) // Initial height for textarea content
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	vp := viewport.New(50, 20)      // For history (initial arbitrary size)
	previewVp := viewport.New(50, 20) // For preview (initial arbitrary size)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &Model{
		messages:    []ai.Message{},
		textarea:    ta,
		viewport:    vp,
		previewVP:   previewVp,
		spinner:     s,
		aiClient:    aiClient,
		ready:       false,
		previewMode: false,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.spinner.Tick)
}

func (m *Model) calculateLayout(width, height int) Layout {
	// Minimum dimensions for overall window
	minOverallWidth := 80  // Reasonable minimum for two columns
	minOverallHeight := 15 // Enough for minimal history, input, and preview

	if width < minOverallWidth {
		width = minOverallWidth
	}
	if height < minOverallHeight {
		height = minOverallHeight
	}

	// Calculate column widths
	leftWidth := int(float64(width) * 0.4) // Left column takes 40%
	// Ensure enough space for borders + minimal content
	minColumnContentWidth := 20
	if leftWidth < minColumnContentWidth+HistoryStyle.GetHorizontalFrameSize() {
		leftWidth = minColumnContentWidth + HistoryStyle.GetHorizontalFrameSize()
	}
	rightWidth := width - leftWidth // Remaining width for right column

	// Adjust rightWidth if it falls below minimum after subtracting leftWidth
	if rightWidth < minColumnContentWidth+PreviewStyle.GetHorizontalFrameSize() {
		rightWidth = minColumnContentWidth + PreviewStyle.GetHorizontalFrameSize()
		leftWidth = width - rightWidth // Re-adjust left to fit if possible
		// If left also shrinks too much, we prioritize right and left might be slightly undersized
		if leftWidth < minColumnContentWidth+HistoryStyle.GetHorizontalFrameSize() {
			leftWidth = minColumnContentWidth + HistoryStyle.GetHorizontalFrameSize()
		}
	}

	// Calculate heights proportionally based on terminal size
	minInputHeight := 8  // Minimum lines needed for input section (title + textarea + help + borders)
	maxInputHeight := 15 // Maximum lines to prevent input from dominating
	minHistoryHeight := 6 // Minimum lines for history section

	// Allocate approximately 25% of height to input, with constraints
	proposedInputHeight := int(float64(height) * 0.25)

	inputHeight := proposedInputHeight
	if inputHeight < minInputHeight {
		inputHeight = minInputHeight
	}
	if inputHeight > maxInputHeight {
		inputHeight = maxInputHeight
	}

	// History gets the remaining space
	historyHeight := height - inputHeight

	// Ensure minimum history height and adjust if necessary
	if historyHeight < minHistoryHeight {
		historyHeight = minHistoryHeight
		inputHeight = height - historyHeight // Adjust input to fit
		// If input becomes too small, we prioritize functionality
		if inputHeight < minInputHeight {
			inputHeight = minInputHeight
			// In very small terminals, sections may overlap or be unusable
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
