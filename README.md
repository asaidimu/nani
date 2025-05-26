# Nani: An AI-Powered Terminal Chat Assistant 🤖

[![Go Reference](https://pkg.go.dev/badge/github.com/asaidimu/nani.svg)](https://pkg.go.dev/github.com/asaidimu/nani)
[![Build Status](https://github.com/asaidimu/nani/workflows/Test%20Workflow/badge.svg)](https://github.com/asaidimu/nani/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Nani is a sleek and efficient terminal-based chat application designed to provide an interactive interface for Google's Gemini AI, featuring real-time markdown rendering and structured AI responses.

---

## ⚡ Quick Links

*   [Overview & Features](#-overview--features)
*   [Installation & Setup](#-installation--setup)
*   [Usage Documentation](#-usage-documentation)
*   [Project Architecture](#-project-architecture)
*   [Development & Contributing](#-development--contributing)
*   [Additional Information](#-additional-information)

---

## ✨ Overview & Features

Nani transforms your terminal into a powerful AI chat environment. Built with Go and the delightful [Charmbracelet](https://charm.sh/) libraries (Bubble Tea, Lipgloss, Glamour), Nani offers a highly interactive and aesthetically pleasing user experience. Its core strength lies in its ability to parse and display rich, structured responses from the Gemini AI, separating the AI's internal thought process, a concise summary, and the detailed content.

This application is particularly adept at handling technical queries, as the integrated Gemini model is pre-configured with a specific persona: an **expert TypeScript developer**. This means Nani is optimized for generating clean, idiomatic TypeScript code, robust interface designs, and providing insightful analysis on modern development practices, making it an invaluable tool for developers seeking intelligent assistance directly from their terminal.

### Key Features:

*   **Interactive Terminal User Interface (TUI)**: A responsive and engaging command-line experience powered by Bubble Tea and Lipgloss.
*   **Google Gemini AI Integration**: Seamless communication with the Gemini AI model for intelligent responses.
*   **Structured AI Responses**: AI output is parsed into distinct `<think>`, `<summary>`, and `<content>` sections, offering clarity and context.
    *   `think`: The AI's detailed reasoning and logical breakdown (visible in chat history).
    *   `summary`: A concise plain-text overview of the response (visible in chat history).
    *   `content`: The complete, detailed answer or solution, often including code blocks (rendered in the real-time preview panel).
*   **Real-time Markdown Preview**: AI-generated markdown content is rendered beautifully in a dedicated preview panel using Glamour.
*   **Responsive Layout**: Adapts dynamically to your terminal window size, ensuring an optimal viewing experience.
*   **Conversation History**: Keeps track of your interactions for continuous context.
*   **Expert AI Persona**: The integrated Gemini model acts as an "Expert TypeScript Developer," providing specialized and high-quality technical assistance.

---

## 🚀 Installation & Setup

To get Nani up and running, you'll need Go installed and a Google Gemini API key.

### Prerequisites

*   **Go**: Version 1.24.3 or higher. You can download it from [golang.org/dl](https://golang.org/dl/).
*   **Google Gemini API Key**: Obtain a key from the [Google AI Studio](https://aistudio.google.com/app/apikey).

### Installation Steps

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/asaidimu/nani.git
    cd nani
    ```
2.  **Install dependencies and build the executable:**
    ```bash
    go mod tidy
    make build
    ```
    This will compile the `nani` executable in the project root directory.

### Configuration

Set your Google Gemini API key as an environment variable named `GEMINI_API_KEY`.

**For Linux/macOS:**

```bash
export GEMINI_API_KEY="YOUR_GEMINI_API_KEY"
```

To make this permanent, add the `export` line to your shell's profile file (e.g., `~/.bashrc`, `~/.zshrc`).

**For Windows (Command Prompt):**

```cmd
set GEMINI_API_KEY="YOUR_GEMINI_API_KEY"
```

**For Windows (PowerShell):**

```powershell
$env:GEMINI_API_KEY="YOUR_GEMINI_API_KEY"
```

### Verification

After setting the API key, you can verify your installation by simply running the `nani` executable:

```bash
./nani
```
You should see the Nani chat interface appear in your terminal. If you encounter an error about `GEMINI_API_KEY` not being set, ensure you've configured the environment variable correctly.

---

## 💡 Usage Documentation

Nani provides a straightforward and intuitive terminal interface for interacting with the Gemini AI.

### Basic Usage

1.  **Start the application:**
    ```bash
    ./nani
    ```
    The terminal UI will launch, showing a "Chat History" panel on the left and a "Preview" panel on the right, along with an input area at the bottom left.

2.  **Type your message:**
    Start typing your query or prompt into the input area. As you type, Nani will prepare to send the message.

3.  **Send your message:**
    Press `Enter` to send your message to the Gemini AI.

4.  **Observe AI response:**
    *   A spinner will appear in the "Chat History" panel indicating that the AI is thinking.
    *   Once the AI responds, the "Chat History" will display "AI: Thinking..." followed by the AI's `summary` and `think` content (combined).
    *   The "Preview" panel will update in real-time with the `content` part of the AI's response, beautifully rendered in markdown.

### Keybindings

*   `Enter`: Send your message to the AI.
*   `Tab`: Toggle between the chat history and the preview panel for focus, or to simply toggle the markdown preview on/off if you prefer.
*   `Q` or `Ctrl+C`: Quit the application.

### Understanding AI Responses

Nani is designed to leverage the structured XML output of the Gemini AI model. When the AI responds, it provides three distinct pieces of information:

*   **`think` (Thought Process)**: This is the AI's internal monologue – its detailed reasoning, assumptions, considerations, and logical steps taken to arrive at the solution. This is displayed in the "Chat History" alongside the summary.
*   **`summary` (Concise Overview)**: A brief, plain-text summary of your request and the AI's response, providing context for past interactions. This is also displayed in the "Chat History."
*   **`content` (Detailed Answer)**: The full, detailed answer or solution, often including code blocks, formatted in Markdown. This is what you'll see rendered in the "Preview" panel.

This separation allows you to quickly grasp the essence of the response (summary), understand the AI's process (think), and review the complete solution (content) simultaneously.

### Example Interaction

```bash
# Start Nani
$ ./nani

# Nani UI appears
# Type your message in the input area:
# You: "Can you define an interface for a User object in TypeScript? It should have id, name, and email."

# Press Enter. Spinner appears.
# After response:

# ----------------------------
# | Chat History           | Preview                    |
# | --------------------   | -------------------------- |
# | You: Can you define    | Welcome to AI Chat         |
# | an interface for a     | Terminal!                  |
# | User object in         |                            |
# | TypeScript? It should  | Features:                  |
# | have id, name, and     | • Real-time markdown       |
# | email.                 |   preview                  |
# |                        | • Responsive layout        |
# | AI: Summary: I have    | • Beautiful terminal UI    |
# | defined a TypeScript   | • AI conversation history  |
# | interface named        |                            |
# | 'User' with id         | Start typing to see your   |
# | (string), name         | message preview here.      |
# | (string), and email    |                            |
# | (string) properties,   | Preview Panel              |
# | along with a brief     | ```typescript              |
# | thought process.       | interface User {           |
# | Thought Process:       |   id: string;              |
# | The request is for a   |   name: string;            |
# | TypeScript interface   |   email: string;           |
# | for a User object      | }                          |
# | with specific fields.  | ```                        |
# | I will create an       |                            |
# | interface and ensure   |                            |
# | appropriate types.     |                            |
# | --------------------   | -------------------------- |
# | Input                  |                            |
# | --------------------   |                            |
# | ┃                      |                            |
# |                        |                            |
# | Enter: Send • Tab:     |                            |
# | Toggle Preview • Q/    |                            |
# | Ctrl+C: Quit           |                            |
# ----------------------------
```

---

## 🏗️ Project Architecture

Nani is a Go application structured for clarity and modularity, primarily leveraging the `charmbracelet` ecosystem for its interactive terminal interface and Google's `genai` SDK for AI integration.

```
nani/
├── main.go               # Application entry point, initializes UI and AI client.
├── go.mod                # Defines module paths and direct dependencies.
├── go.sum                # Checksums for module dependencies.
├── Makefile              # Standard build, test, and clean commands.
└── pkg/                  # Contains core application logic packages.
    ├── ai/               # Handles all AI-related logic.
    │   ├── gemini.go     # Implements the AIClient interface using Google Gemini API.
    │   ├── types.go      # Defines AI message structure, AIClient interface, and structured AI Response.
    │   └── utils.go      # Utility functions for parsing structured AI responses (XML).
    └── ui/               # Manages the Terminal User Interface (TUI).
        ├── model.go      # Defines the Bubble Tea model (application state) and layout calculations.
        ├── styles.go     # Contains Lipgloss styles for all UI elements.
        ├── update.go     # Implements Bubble Tea's Update method, handling user input and AI responses.
        └── view.go       # Implements Bubble Tea's View method, rendering the TUI.
```

### Core Components

*   **`main.go`**: The application's entry point. It sets up the Gemini AI client, initializes the TUI model, and starts the Bubble Tea program. It also handles environment variable checks for the API key.
*   **`pkg/ai`**:
    *   **`AIClient` Interface**: Defines the contract for any AI service integration, allowing for potential future AI model swaps.
    *   **`GeminiAIClient`**: The concrete implementation of `AIClient` for Google's Gemini API. This is where the specific AI persona (Expert TypeScript Developer) and the mandatory XML response structure are embedded in the system prompt.
    *   **Response Struct**: Dictates the expected structured format (`<response>`, `<think>`, `<summary>`, `<content>`) from the AI.
    *   **XML Parsing**: Utilities within this package ensure that the AI's raw text response is correctly parsed into the structured `Response` object.
*   **`pkg/ui`**: This package encapsulates all terminal UI logic using the `charmbracelet` libraries.
    *   **`Model`**: Holds the entire state of the TUI, including messages, text area, viewports, loading status, and layout dimensions. It also manages the responsive resizing of UI elements.
    *   **`Update`**: The heart of the Bubble Tea application, processing user inputs (key presses) and internal messages (AI responses, window resize events) to update the model state. It initiates AI requests in a non-blocking manner.
    *   **`View`**: Renders the current state of the model to the terminal, arranging chat history, input area, and the markdown preview panel.
    *   **`Styles`**: Centralized definitions for all visual styles, colors, and borders using Lipgloss, ensuring a consistent and appealing aesthetic.

### Data Flow

1.  **User Input**: The user types a message into the `textarea` and presses `Enter`.
2.  **UI `Update`**: The `Update` function in `pkg/ui` receives the `tea.KeyMsg`. It updates the `Model` to include the user's message in history, clears the `textarea`, sets a `loading` flag, and triggers a `sendToAI` command.
3.  **AI Request (Goroutine)**: The `sendToAI` command executes as a `tea.Cmd`, which runs in a separate goroutine. It calls `m.aiClient.SendMessage` with the user's message and current conversation history.
4.  **Gemini Interaction**: The `GeminiAIClient` sends the message to the Google Gemini API.
5.  **Structured Response**: The Gemini API responds. `GeminiAIClient` then uses `parseAIResponse` to transform the raw text into the structured `Response` (`think`, `summary`, `content`) object.
6.  **AI Response Message**: The `AIResponseMsg` (containing the structured AI response) is sent back to the main Bubble Tea event loop.
7.  **UI `Update` (AI Response)**: The `Update` function processes `AIResponseMsg`. It updates the `messages` history to include the AI's `summary` and `think` content, and updates the `previewVP` with the markdown-rendered `content`. The `loading` flag is reset.
8.  **UI `View`**: After each `Update`, the `View` function is called to re-render the entire TUI based on the new `Model` state, displaying the updated chat history and preview.

---

## 🛠️ Development & Contributing

Contributions are welcome! Whether it's a bug report, a new feature, or an improvement to the documentation, your input is valuable.

### Development Setup

1.  **Fork** the `asaidimu/nani` repository on GitHub.
2.  **Clone** your forked repository:
    ```bash
    git clone https://github.com/YOUR_USERNAME/nani.git
    cd nani
    ```
3.  **Install Go modules**:
    ```bash
    go mod tidy
    ```
4.  **Build the project**:
    ```bash
    make build
    ```
    This will create the `nani` executable.

### Scripts

The project includes a `Makefile` for common development tasks:

*   `make all` (default): Runs `make build`.
*   `make build`: Compiles the `nani` executable.
*   `make test`: Runs all tests in the project.
*   `make clean`: Removes the compiled `nani` executable.

### Testing

To run the test suite, simply execute:

```bash
make test
```
Ensure all tests pass before submitting a pull request.

### Contributing Guidelines

We follow a [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification for commit messages. This helps in generating changelogs and automating semantic versioning.

*   **fix:** a commit that fixes a bug (corresponds to PATCH in SemVer)
*   **feat:** a commit that adds new functionality (corresponds to MINOR in SemVer)
*   **feat!:** or **fix!:** or **refactor!:** etc., a commit with a footer `BREAKING CHANGE:` introduces a breaking API change (corresponds to MAJOR in SemVer)

1.  **Fork** the repository and create your feature branch from `main` (`git checkout -b feature/your-feature-name`).
2.  **Commit** your changes using the Conventional Commits format (e.g., `feat: add new CLI command`).
3.  **Push** your branch (`git push origin feature/your-feature-name`).
4.  **Open a Pull Request** against the `main` branch of the upstream repository.
5.  Ensure your code adheres to Go formatting standards (`go fmt ./...`) and passes all tests.

### Issue Reporting

Encounter a bug or have a feature idea? Please open an issue on the [GitHub Issue Tracker](https://github.com/asaidimu/nani/issues). Provide as much detail as possible, including steps to reproduce, expected vs. actual behavior, and your environment setup.

---

## 📚 Additional Information

### Troubleshooting

*   **`Error: GEMINI_API_KEY environment variable not set`**: Ensure you have set the `GEMINI_API_KEY` environment variable correctly before running `nani`. Double-check for typos and that it's accessible in your terminal session.
*   **"Failed to create Gemini client" / API errors**: Verify your `GEMINI_API_KEY` is valid and has the necessary permissions for the Gemini API. Check your internet connection.
*   **UI rendering issues**: Ensure your terminal emulator supports 256 colors and Unicode characters. Older terminals might have display glitches. Try resizing your terminal window.

### Changelog / Roadmap

For upcoming features and changes, please refer to the project's [GitHub Releases](https://github.com/asaidimu/nani/releases) and the [Issues](https://github.com/asaidimu/nani/issues) for planned work.

### License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

### Acknowledgments

Nani is built upon the incredible work of the following open-source projects:

*   [Charmbracelet Bubble Tea](https://github.com/charmbracelet/bubbletea): A powerful framework for building TUIs.
*   [Charmbracelet Lipgloss](https://github.com/charmbracelet/lipgloss): A library for styling terminal output.
*   [Charmbracelet Glamour](https://github.com/charmbracelet/glamour): A markdown renderer for the terminal.
*   [Google for Go SDK for Gemini](https://github.com/google/generative-ai-go): The official Go client for the Gemini API.