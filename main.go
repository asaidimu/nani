package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/asaidimu/nani/pkg/ai"
	"github.com/asaidimu/nani/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: GEMINI_API_KEY environment variable not set")
		os.Exit(1)
	}

	project :=  filepath.Join(".")
	workspace, err := ai.NewWorkspace(project)
	if err != nil {
		fmt.Printf("Error creating workspace: %v\n", err)
		os.Exit(1)
	}

	err = workspace.Init("nani", "saidimu", "https://github.com/asaidimu/nani.git")
	if err != nil {
		fmt.Printf("Error initializing workspace: %v\n", err)
		os.Exit(1)
	}

	aiClient, err := ai.NewGeminiAIClient(apiKey, workspace)
	if err != nil {
		fmt.Printf("Error initializing Gemini client: %v\n", err)
		os.Exit(1)
	}

	m := ui.New(aiClient)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
