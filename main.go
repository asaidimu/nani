// main.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/asaidimu/nani/pkg/ai"
	"github.com/asaidimu/nani/pkg/ui" // Import the new ui package
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: GEMINI_API_KEY environment variable not set")
		os.Exit(1)
	}
	aiClient, err := ai.NewGeminiAIClient(apiKey)
	if err != nil {
		fmt.Printf("Error initializing Gemini client: %v\n", err)
		os.Exit(1)
	}
	ctx := context.Background();
	_, err = aiClient.StartSession(ctx)
	if err != nil {
		fmt.Printf("Error initializing Gemini client: %v\n", err)
		os.Exit(1)
	}

	m := ui.New(aiClient) // Use the New function from the ui package

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
