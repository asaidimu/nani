// Package ai provides a structured file-based persistence layer for managing
// AI application workspace configurations, session data, defined roles,
// and user preferences. It handles the creation and maintenance of a
// .AIWorkspace directory and uses an in-memory indexing system for efficient
// retrieval of artifact summaries.
package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SessionSummary provides a lightweight summary of an archived session.
// It is used primarily for listing available sessions without loading their
// entire content (like chat history or source code lists).
type SessionSummary struct {
	ID          string    `json:"id"`          // Unique identifier for the session.
	Label       string    `json:"label"`       // A human-readable label for the session.
	RoleName    string    `json:"roleName"`    // The name of the AI role used in this session.
	CreatedAt   time.Time `json:"createdAt"`   // Timestamp when the session was created.
	LastUpdated time.Time `json:"lastUpdated"` // Timestamp when the session was last updated.
}

// RoleSummary provides a lightweight summary of an AI role.
// It is used for listing available roles without loading the full persona string.
type RoleSummary struct {
	Name        string `json:"name"`        // Unique name of the role (e.g., "documenter").
	Label       string `json:"label"`       // Human-readable label for the role (e.g., "Code Documenter").
	Description string `json:"description"` // A brief description of the role's purpose.
}

// PreferenceSummary provides a lightweight summary of a user preference.
// It is used for listing preferences, including a snippet of their content.
type PreferenceSummary struct {
	ID             string    `json:"id"`                       // Unique identifier for the preference.
	Timestamp      time.Time `json:"timestamp"`                // Timestamp when the preference was created or last updated.
	ContentSnippet string    `json:"contentSnippet,omitempty"` // A truncated snippet of the preference's content.
}

// ArtifactIndexes groups all artifact indexes together within the workspace context.
// This provides a centralized and organized way to quickly access summaries of
// various stored data types without reading full files from disk for every query.
type ArtifactIndexes struct {
	ArchivedSessions map[string]SessionSummary   `json:"sessions"`   // Index of archived sessions, keyed by session ID.
	RolesIndex       map[string]RoleSummary      `json:"roles"`      // Index of roles, keyed by role name.
	PreferencesIndex map[string]PreferenceSummary `json:"preferences"`// Index of preferences, keyed by preference ID.
}

// Context represents the overall workspace configuration.
// It is stored in `context.json` and includes global settings, project metadata,
// and in-memory indexes of various artifacts for quick lookup.
type Context struct {
	Workspace string          `json:"workspace"` // A unique ID for the workspace itself.
	Settings  Settings        `json:"settings"`  // Workspace-wide settings.
	Project   Project         `json:"project"`   // Project-specific metadata.
	Indexes   ArtifactIndexes `json:"indexes"`   // Nested indexes for better organization and quick lookup.
}

// Settings holds workspace-wide configuration settings.
type Settings struct {
	DefaultLanguage string `json:"defaultLanguage"` // The default language setting for the AI.
	DefaultRole     string `json:"defaultRole"`     // The name of the default AI role to use.
	SystemPrompt    string `json:"systemPrompt"`    // A global system prompt applied to all AI interactions.
}

// Project holds metadata specific to the AI project associated with the workspace.
type Project struct {
	Name       string `json:"name"`       // The name of the project.
	Owner      string `json:"owner"`      // The owner or creator of the project.
	Repository string `json:"repository"` // URL or identifier of the project's source code repository.
}

// Session represents an active or archived interaction session with the AI.
// Active sessions are stored in `session.json`, while archived sessions are
// moved to `sessions/<id>.json`.
type Session struct {
	ID       string   `json:"id"`       // Unique identifier for this session.
	Label    string   `json:"label"`    // A descriptive label for the session.
	Role     Role     `json:"role"`     // The full AI role configuration for this session.
	Sources  []string `json:"sources"`  // A list of file paths that are relevant to this session.
	Chat     []Chat   `json:"chat"`     // A chronological list of user-AI interactions.
	Metadata Metadata `json:"metadata"` // Internal session management data.
}

// MarshalJSON customizes Session JSON serialization.
// It ensures that only the `Role.Name` is saved to JSON for the `Role` field,
// rather than the entire `Role` struct, keeping the session file compact.
func (s Session) MarshalJSON() ([]byte, error) {
	type Alias Session // Create an alias to prevent infinite recursion
	// When marshaling, 's' is a value receiver. To create a pointer to Alias from 's',
	// we need to take the address of 's' first.
	aux := (*Alias)(&s) // Correct: Convert pointer to s to pointer to Alias

	return json.Marshal(&struct {
		*Alias
		Role string `json:"role"` // This field will hold the Role.Name for JSON serialization.
	}{
		Alias: aux, // Assign the *Alias pointer
		Role:  s.Role.Name, // Store only the role's name.
	})
}

// UnmarshalJSON customizes Session JSON deserialization.
// It populates the `Role` field by unmarshaling only the role's name initially.
// The full `Role` struct data (Persona, Description, etc.) is subsequently loaded
// by `Workspace.loadSession` using this role name, ensuring the `Role` is complete in memory.
func (s *Session) UnmarshalJSON(data []byte) error {
	type Alias Session // Create an alias to prevent infinite recursion
	aux := &struct {
		*Alias
		RoleName string `json:"role"` // Temporary field to unmarshal the role's name from JSON.
	}{
		Alias: (*Alias)(s), // This is correct because 's' is already a pointer (*Session).
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Role = Role{Name: aux.RoleName} // Populate only the Name; full Role struct is loaded later.
	return nil
}

// Chat represents a single user-AI interaction within a session.
type Chat struct {
	ID       string        `json:"id"`       // Unique identifier for this chat interaction.
	Message  SavedMessage  `json:"message"`  // The user's input message.
	Response SavedResponse `json:"response"` // The AI's response to the message.
}

// SavedMessage is a user's prompt or input, stored persistently.
type SavedMessage struct {
	Content   string    `json:"content"`   // The textual content of the user's message.
	Timestamp time.Time `json:"timestamp"` // The timestamp when the message was created.
}

// SavedResponse is the AI's reply to a user's message, stored persistently.
type SavedResponse struct {
	Content   string    `json:"content"`   // The textual content of the AI's response.
	Timestamp time.Time `json:"timestamp"` // The timestamp when the response was generated.
}

// Metadata holds internal management data for a session, useful for tracking
// its lifecycle and characteristics.
type Metadata struct {
	CreatedAt       time.Time `json:"createdAt"`       // Timestamp when the session was originally created.
	Priority        string    `json:"priority"`        // Indication of session importance (e.g., "low", "medium", "high").
	SessionDuration string    `json:"sessionDuration"` // Expected or actual duration of the session in seconds (as string).
	LastUpdated     time.Time `json:"lastUpdated"`     // Timestamp of the last modification to the session.
	ArchiveAfter    time.Time `json:"archiveAfter"`    // Timestamp after which the session is eligible for archiving.
}

// Preference represents a user-defined AI prompt tweak or instruction.
// Preferences are stored as individual JSON files in the `preferences/` directory.
type Preference struct {
	ID        string    `json:"id"`        // Unique identifier for the preference.
	Content   string    `json:"content"`   // The detailed textual content of the preference.
	Timestamp time.Time `json:"timestamp"` // The timestamp when the preference was created or last updated.
}

// Role represents an AI persona or configuration.
// Roles define how the AI should behave and are stored as individual JSON files
// in the `roles/` directory.
type Role struct {
	Name        string `json:"name"`        // Unique name of the role (e.g., "documenter").
	Label       string `json:"label"`       // Human-readable label for the role (e.g., "Code Documenter").
	Persona     string `json:"persona"`     // The detailed prompt string that defines the AI's personality/instructions.
	Description string `json:"description"` // A brief description of the role's purpose.
}

// Workspace manages the `.AIWorkspace` directory, which serves as the root
// for all persistent data for an AI application. It provides methods for
// initializing the workspace, managing sessions, roles, and preferences.
type Workspace struct {
	RootDir string  // The root directory where `.AIWorkspace` is located.
	Context Context // The in-memory representation of the workspace's context.
}

// NewWorkspace creates a new Workspace instance.
// It initializes the `.AIWorkspace` directory and its required subdirectories
// (`preferences`, `sessions`, `roles`, `logs`) if they don’t already exist.
// This function primarily handles the physical setup of the workspace directory structure.
func NewWorkspace(rootDir string) (*Workspace, error) {
	aiDir := filepath.Join(rootDir, ".AIWorkspace")

	// Check if .AIWorkspace exists, create if not
	if _, err := os.Stat(aiDir); os.IsNotExist(err) {
		if err := os.MkdirAll(aiDir, 0755); err != nil { // 0755: owner rwx, group rx, others rx
			return nil, fmt.Errorf("failed to create workspace directory %s: %w", aiDir, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check workspace directory %s: %w", aiDir, err)
	}

	// Ensure subdirectories exist
	for _, dir := range []string{"preferences", "sessions", "roles", "logs"} {
		subDir := filepath.Join(aiDir, dir)
		if _, err := os.Stat(subDir); os.IsNotExist(err) {
			if err := os.MkdirAll(subDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create %s directory: %w", dir, err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("failed to check %s directory: %w", dir, err)
		}
	}

	return &Workspace{
		RootDir: aiDir,
	}, nil
}

// Init initializes a new workspace project, or loads an existing one.
// It checks for `context.json`, creates a default one if it doesn't exist,
// or loads the existing one. It ensures default roles are present and
// rebuilds all in-memory artifact indexes to synchronize with disk.
// This method is typically called once at application startup.
func (w *Workspace) Init(projectName, owner, repo string) error {
	contextPath := filepath.Join(w.RootDir, "context.json")

	// Flag to track if a new context was created
	newContextCreated := false

	// Check if context.json exists
	if _, err := os.Stat(contextPath); os.IsNotExist(err) {
		// Create default context if not found
		context := Context{
			Workspace: uuid.New().String(),
			Settings: Settings{
				DefaultLanguage: "en",
				DefaultRole:     "documenter",
				SystemPrompt:    "You are a general-purpose AI assistant. Provide concise and helpful responses.", // Default system prompt
			},
			Project: Project{Name: projectName, Owner: owner, Repository: repo},
			Indexes: ArtifactIndexes{ // Initialize nested struct for indexes
				ArchivedSessions: make(map[string]SessionSummary),
				RolesIndex:       make(map[string]RoleSummary),
				PreferencesIndex: make(map[string]PreferenceSummary),
			},
		}
		if err := w.saveContext(context); err != nil {
			return fmt.Errorf("failed to save context: %w", err)
		}
		w.Context = context // Set the workspace's context
		newContextCreated = true
	} else if err != nil {
		return fmt.Errorf("failed to check context file %s: %w", contextPath, err)
	} else {
		// Load existing context
		if err := w.loadContext(); err != nil {
			return fmt.Errorf("failed to load context: %w", err)
		}
	}

	// Backward compatibility: If SystemPrompt is empty in an existing context, set a default.
	// This handles cases where old context.json files don't have this field.
	if w.Context.Settings.SystemPrompt == "" {
		w.Context.Settings.SystemPrompt = "You are a general-purpose AI assistant. Provide concise and helpful responses."
		// Save context immediately if system prompt was missing and set, to persist the default.
		if err := w.saveContext(w.Context); err != nil {
			return fmt.Errorf("failed to update context with default system prompt: %w", err)
		}
	}


	// Ensure index maps are initialized if loaded context had nil maps (e.g., from old schema or if 'Indexes' struct was nil)
	if w.Context.Indexes.ArchivedSessions == nil {
		w.Context.Indexes.ArchivedSessions = make(map[string]SessionSummary)
	}
	if w.Context.Indexes.RolesIndex == nil {
		w.Context.Indexes.RolesIndex = make(map[string]RoleSummary)
	}
	if w.Context.Indexes.PreferencesIndex == nil {
		w.Context.Indexes.PreferencesIndex = make(map[string]PreferenceSummary)
	}

	// Rebuild/Reconcile indexes (important for new workspaces or schema migrations from old schema)
	// Only rebuild if a new context wasn't just created (as it would be empty anyway)
	// or if we are loading an existing context which might be out of sync.
	if !newContextCreated {
		if err := w.rebuildIndexes(); err != nil {
			return fmt.Errorf("failed to rebuild indexes: %w", err)
		}
	}


	// Create default documenter role if its file doesn't exist.
	// This will also add it to the index via saveRole.
	rolePath := filepath.Join(w.RootDir, "roles", "documenter.json")
	if _, err := os.Stat(rolePath); os.IsNotExist(err) {
		role := Role{
			Name:        "documenter",
			Label:       "Code Documenter",
			Persona:     "You are a meticulous technical writer who creates clear, detailed markdown documentation with a high level of verbosity, including examples where appropriate, and adheres to user-specified preferences.",
			Description: "Generates detailed documentation for code files, tailored to user preferences in markdown format.",
		}
		if err := w.saveRole(role); err != nil { // saveRole will update the index
			return fmt.Errorf("failed to save default role: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check documenter role file %s: %w", rolePath, err)
	}

	return w.logAction("Initialized workspace")
}

// rebuildIndexes scans the file system directories for sessions, roles, and preferences
// and rebuilds the in-memory indexes within the Workspace's Context.
// This is an internal helper function called by `Init()` and `RefreshIndexes()`.
func (w *Workspace) rebuildIndexes() error {
	// Re-initialize all index maps to ensure a clean rebuild
	w.Context.Indexes.ArchivedSessions = make(map[string]SessionSummary)
	w.Context.Indexes.RolesIndex = make(map[string]RoleSummary)
	w.Context.Indexes.PreferencesIndex = make(map[string]PreferenceSummary)

	// Rebuild session index
	sessionsDir := filepath.Join(w.RootDir, "sessions")
	files, err := os.ReadDir(sessionsDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read sessions directory for rebuilding index: %w", err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			sessionPath := filepath.Join(sessionsDir, file.Name())
			data, err := os.ReadFile(sessionPath)
			if err != nil {
				w.logAction(fmt.Sprintf("Warning: Could not read archived session file '%s' during index rebuild: %v\n", sessionPath, err))
				continue // Continue processing other files
			}
			// Use a temporary anonymous struct for unmarshaling just the summary parts
			temp := struct {
				ID       string   `json:"id"`
				Label    string   `json:"label"`
				Role     string   `json:"role"` // Unmarshal role name from JSON
				Metadata Metadata `json:"metadata"`
			}{}
			if err := json.Unmarshal(data, &temp); err != nil {
				w.logAction(fmt.Sprintf("Warning: Could not parse archived session summary from '%s' during index rebuild: %v\n", sessionPath, err))
				continue // Continue processing other files
			}

			// Create a SessionSummary from the parsed data
			w.Context.Indexes.ArchivedSessions[temp.ID] = SessionSummary{
				ID:        temp.ID,
				Label:     temp.Label,
				RoleName:  temp.Role, // Use the unmarshaled role name
				CreatedAt: temp.Metadata.CreatedAt,
				LastUpdated: temp.Metadata.LastUpdated,
			}
		}
	}

	// Rebuild roles index
	rolesDir := filepath.Join(w.RootDir, "roles")
	files, err = os.ReadDir(rolesDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read roles directory for rebuilding index: %w", err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			rolePath := filepath.Join(rolesDir, file.Name())
			data, err := os.ReadFile(rolePath)
			if err != nil {
				w.logAction(fmt.Sprintf("Warning: Could not read role file '%s' during index rebuild: %v\n", rolePath, err))
				continue
			}
			var r Role
			if err := json.Unmarshal(data, &r); err != nil {
				w.logAction(fmt.Sprintf("Warning: Could not parse role from '%s' during index rebuild: %v\n", rolePath, err))
				continue
			}
			w.Context.Indexes.RolesIndex[r.Name] = RoleSummary{
				Name:        r.Name,
				Label:       r.Label,
				Description: r.Description,
			}
		}
	}

	// Rebuild preferences index
	preferencesDir := filepath.Join(w.RootDir, "preferences")
	files, err = os.ReadDir(preferencesDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read preferences directory for rebuilding index: %w", err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			prefPath := filepath.Join(preferencesDir, file.Name())
			data, err := os.ReadFile(prefPath)
			if err != nil {
				w.logAction(fmt.Sprintf("Warning: Could not read preference file '%s' during index rebuild: %v\n", prefPath, err))
				continue
			}
			var p Preference
			if err := json.Unmarshal(data, &p); err != nil {
				w.logAction(fmt.Sprintf("Warning: Could not parse preference from '%s' during index rebuild: %v\n", prefPath, err))
				continue
				}
			snippet := p.Content
			if len(snippet) > 100 { // Limit snippet length for snippet
				snippet = snippet[:100] + "..."
			}
			w.Context.Indexes.PreferencesIndex[p.ID] = PreferenceSummary{
				ID:             p.ID,
				Timestamp:      p.Timestamp,
				ContentSnippet: snippet,
			}
		}
	}

	// After rebuilding, save the context to persist the new indexes
	return w.saveContext(w.Context)
}

// RefreshIndexes explicitly triggers a re-scan of the artifact directories and rebuilds the in-memory indexes.
// This method can be called by the user of the package if manual changes to artifact files (roles, preferences, archived sessions)
// are suspected or have occurred outside of the package's direct API calls, to synchronize the in-memory state.
// It performs a synchronous operation. For non-blocking behavior, call it within a goroutine from your application.
func (w *Workspace) RefreshIndexes() error {
	w.logAction("Refreshing workspace indexes initiated.")
	if err := w.rebuildIndexes(); err != nil {
		return fmt.Errorf("failed to refresh indexes: %w", err)
	}
	return w.logAction("Workspace indexes refreshed successfully.")
}

// GetSession retrieves the currently active session. If no active session is found,
// a new session is automatically created using the provided `defaultLabel` and `defaultRoleName`.
//
// The `defaultLabel` will be used as the label for the new session if one is created.
// The `defaultRoleName` allows specifying a role other than the `DefaultRole`
// configured in `Context.Settings` for the newly created session. If `defaultRoleName` is empty,
// or if the specified role cannot be found, the `DefaultRole` from settings will be used.
//
// This method encapsulates the common pattern of ensuring an active session is always available.
func (w *Workspace) GetSession(defaultLabel string, defaultRoleName string) (*Session, error) {
	session, err := w.GetActiveSession()
	if err != nil {
		return nil, fmt.Errorf("error checking for active session: %w", err)
	}

	if session == nil {
		// StartSession handles default role fallback if defaultRoleName is empty or invalid
		newSession, createErr := w.StartSession(defaultLabel, defaultRoleName)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create new session: %w", createErr)
		}
		return newSession, nil
	}

	return session, nil
}

// StartSession begins a new interaction session.
// If an active session (`session.json`) already exists, it is first archived
// using `EndSession()` before a new session is created.
//
// The `desiredRoleName` parameter allows specifying a role other than the `DefaultRole`
// configured in `Context.Settings`. If `desiredRoleName` is empty, or if the specified
// role cannot be found, the `DefaultRole` will be used.
//
// The new session is initialized with a unique ID, a human-readable label,
// the determined role, and current metadata. The active session data is saved to `session.json`.
func (w *Workspace) StartSession(label string, desiredRoleName string) (*Session, error) {
	sessionPath := filepath.Join(w.RootDir, "session.json")

	// Archive existing session if present
	if _, err := os.Stat(sessionPath); err == nil {
		if err := w.EndSession(); err != nil { // EndSession will update the index
			return nil, fmt.Errorf("failed to archive existing session: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// Handle other errors when checking session.json
		return nil, fmt.Errorf("failed to check for active session: %w", err)
	}

	// Determine which role to use
	roleToUse := w.Context.Settings.DefaultRole
	if desiredRoleName != "" {
		// Check if the desired role exists in our index
		_, found := w.Context.Indexes.RolesIndex[desiredRoleName]
		if found {
			roleToUse = desiredRoleName
		} else {
			// Log a warning if the desired role wasn't found and fallback to default
			w.logAction(fmt.Sprintf("Warning: Desired role '%s' not found. Falling back to default role '%s'.\n",
				desiredRoleName, w.Context.Settings.DefaultRole))
		}
	}

	// Load the determined role's full data from disk
	role, err := w.loadRole(roleToUse)
	if err != nil {
		// This error should ideally not happen if roleToUse is from default or found in index.
		return nil, fmt.Errorf("failed to load role '%s' for new session: %w", roleToUse, err)
	}

	// Create new session
	now := time.Now()
	session := &Session{
		ID:      uuid.New().String(),
		Label:   label,
		Role:    role,
		Sources: []string{},
		Chat:    []Chat{},
		Metadata: Metadata{
			CreatedAt:       now,
			Priority:        "medium",           // Example default priority
			SessionDuration: "3600",             // Example default: 1 hour in seconds as string
			LastUpdated:     now,
			ArchiveAfter:    now.Add(7 * 24 * time.Hour), // Automatically archive after 7 days
		},
	}
	if err := w.saveSession(*session); err != nil {
		return nil, fmt.Errorf("failed to save new session: %w", err)
	}

	if err := w.logAction(fmt.Sprintf("Started session %s with label '%s' and role '%s'", session.ID, session.Label, role.Name)); err != nil {
		return nil, fmt.Errorf("failed to log session start: %w", err)
	}

	return session, nil
}

// EndSession archives the current active session.
// The `session.json` file is moved to the `sessions/` subdirectory (named `sessions/<id>.json`),
// and its summary is added to the `ArchivedSessions` index in the `Context`.
// The `session.json` file is then removed. If no active session exists, the method does nothing.
func (w *Workspace) EndSession() error {
	sessionPath := filepath.Join(w.RootDir, "session.json")
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		return nil // No active session to archive, gracefully exit
	}

	session, err := w.loadSession(); // loadSession handles Role hydration
	if err != nil {
		return fmt.Errorf("failed to load session for archiving: %w", err)
	}

	// Save to sessions/<id>.json
	archivePath := filepath.Join(w.RootDir, "sessions", fmt.Sprintf("%s.json", session.ID))
	if err := w.writeJSON(archivePath, session); err != nil {
		return fmt.Errorf("failed to archive session %s: %w", session.ID, err)
	}

	// Remove session.json
	if err := os.Remove(sessionPath); err != nil {
		return fmt.Errorf("failed to remove active session file %s after archiving: %w", sessionPath, err)
	}

	// Add to archived sessions index
	w.Context.Indexes.ArchivedSessions[session.ID] = SessionSummary{
		ID:        session.ID,
		Label:     session.Label,
		RoleName:  session.Role.Name,
		CreatedAt: session.Metadata.CreatedAt,
		LastUpdated: session.Metadata.LastUpdated,
	}
	if err := w.saveContext(w.Context); err != nil {
		return fmt.Errorf("failed to update context after archiving session: %w", err)
	}

	return w.logAction(fmt.Sprintf("Archived session %s", session.ID))
}

// AddSource adds a source file path to the `Sources` list of the current active session.
// It validates that the source file exists and ensures no duplicate paths are added.
// The session's `LastUpdated` timestamp is updated, and the session is saved back to disk.
func (w *Workspace) AddSource(sourcePath string) error {
	session, err := w.loadSession(); // loadSession handles Role hydration
	if err != nil {
		return fmt.Errorf("failed to load session to add source: %w", err)
	}

	// Validate source path (basic check for existence)
	// Consider making this path absolute or relative to RootDir for consistency if not already.
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source file %s does not exist: %w", sourcePath, err)
	} else if err != nil {
		return fmt.Errorf("failed to stat source file %s: %w", sourcePath, err)
	}

	// Add source if not already present
	for _, src := range session.Sources {
		if src == sourcePath {
			return nil // Source already added, no action needed
		}
	}
	session.Sources = append(session.Sources, sourcePath)
	session.Metadata.LastUpdated = time.Now()

	if err := w.saveSession(*session); err != nil {
		return fmt.Errorf("failed to save session after adding source %s: %w", sourcePath, err)
	}

	return w.logAction(fmt.Sprintf("Added source %s to session %s", sourcePath, session.ID))
}

// AddInteraction adds a user-AI interaction to the `Chat` history of the current active session.
// A new `Chat` entry is created with the provided user prompt and AI response,
// and the session's `LastUpdated` timestamp is updated. The session is saved back to disk.
func (w *Workspace) AddInteraction(userPrompt, aiResponse string) error {
	session, err := w.loadSession(); // loadSession handles Role hydration
	if err != nil {
		return fmt.Errorf("failed to load session to add interaction: %w", err)
	}

	// Create new chat entry
	now := time.Now()
	chat := Chat{
		ID: uuid.New().String(),
		Message: SavedMessage{
			Content:   userPrompt,
			Timestamp: now,
		},
		Response: SavedResponse{
			Content:   aiResponse,
			Timestamp: now.Add(1 * time.Second), // Slight offset for response timestamp
		},
	}

	// Append chat and update metadata
	session.Chat = append(session.Chat, chat)
	session.Metadata.LastUpdated = now

	if err := w.saveSession(*session); err != nil {
		return fmt.Errorf("failed to save session after adding interaction: %w", err)
	}

	return w.logAction(fmt.Sprintf("Added interaction (chat ID: %s) to session %s", chat.ID, session.ID))
}

// SwitchRole changes the AI role for the current active session.
// It loads the new role configuration from disk, updates the session's `Role` field
// and `LastUpdated` timestamp, and saves the session back to disk.
func (w *Workspace) SwitchRole(roleName string) error {
	session, err := w.loadSession(); // loadSession handles Role hydration
	if err != nil {
		return fmt.Errorf("failed to load session to switch role: %w", err)
	}

	// Load new role
	role, err := w.loadRole(roleName)
	if err != nil {
		return fmt.Errorf("failed to load role %s for switching: %w", roleName, err)
	}

	// Update session role and metadata
	session.Role = role
	session.Metadata.LastUpdated = time.Now()

	if err := w.saveSession(*session); err != nil {
		return fmt.Errorf("failed to save session after switching to role %s: %w", roleName, err)
	}

	return w.logAction(fmt.Sprintf("Switched session %s to role %s", session.ID, roleName))
}

// GetActiveSession loads and returns the current active session.
// It returns a pointer to the `Session` struct if `session.json` exists and can be parsed.
// If no active session is found (i.e., `session.json` does not exist), it returns `nil, nil`.
// An error is returned if `session.json` exists but cannot be read or parsed.
func (w *Workspace) GetActiveSession() (*Session, error) {
	session, err := w.loadSession()
	if err != nil {
		// Specifically check for the "no active session found" error by message content
		if os.IsNotExist(err) || strings.Contains(err.Error(), "no active session found") {
			return nil, nil // No active session, not an error state for this public API
		}
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}
	return session, nil
}

// ResumeArchivedSession moves an archived session back to the active `session.json` state.
// If an active session currently exists, it is first archived using `EndSession()`.
// The specified archived session file is read, parsed, made the new active session,
// its summary is removed from the `ArchivedSessions` index, and the original archived file is optionally removed.
func (w *Workspace) ResumeArchivedSession(sessionID string) (*Session, error) {
	// First, archive any currently active session to ensure a clean state
	if err := w.EndSession(); err != nil {
		return nil, fmt.Errorf("failed to archive current session before resuming archived one: %w", err)
	}

	archivePath := filepath.Join(w.RootDir, "sessions", fmt.Sprintf("%s.json", sessionID))

	// Check if the archived session file exists
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("archived session with ID '%s' not found at '%s': %w", sessionID, archivePath, err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to check archived session file '%s': %w", archivePath, err)
	}

	// Load the archived session data
	data, err := os.ReadFile(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archived session file '%s': %w", archivePath, err)
	}

	var session Session
	// Unmarshal the archived session data (Session.UnmarshalJSON will only populate Role.Name)
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse archived session data from '%s': %w", archivePath, err)
	}

	// Load the full role data for the session's role name
	role, err := w.loadRole(session.Role.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to load full role data for archived session '%s' (role name: '%s'): %w", sessionID, session.Role.Name, err)
	}
	session.Role = role // Assign the fully loaded role to the session

	// Save the loaded archived session as the new active session (session.json)
	if err := w.saveSession(session); err != nil {
		return nil, fmt.Errorf("failed to save archived session '%s' as active session: %w", sessionID, err)
	}

	// Remove from archived sessions index in Context
	delete(w.Context.Indexes.ArchivedSessions, session.ID)
	if err := w.saveContext(w.Context); err != nil {
		return nil, fmt.Errorf("failed to update context after resuming session: %w", err)
	}

	// Optionally, remove the original archived file if the intent is to "move" it, not copy.
	if err := os.Remove(archivePath); err != nil {
		// Log this as a warning, but don't fail the entire resume operation as the active session is now set.
		w.logAction(fmt.Sprintf("Warning: Failed to remove original archived session file '%s' after resuming: %v\n", archivePath, err))
	}

	// Log the successful resumption of the session
	if err := w.logAction(fmt.Sprintf("Resumed archived session %s", sessionID)); err != nil {
		return nil, fmt.Errorf("failed to log session resume for ID '%s': %w", sessionID, err)
	}

	return &session, nil
}

// ListArchivedSessions returns a slice of all archived session summaries.
// This data is retrieved directly from the in-memory `ArchivedSessions` index in the `Context`,
// making it a very efficient operation as it avoids reading individual session files from disk.
func (w *Workspace) ListArchivedSessions() ([]SessionSummary, error) {
	// Convert map values to slice
	sessions := make([]SessionSummary, 0, len(w.Context.Indexes.ArchivedSessions))
	for _, s := range w.Context.Indexes.ArchivedSessions {
		sessions = append(sessions, s)
	}
	return sessions, nil
}

// ListRoles returns a slice of all role summaries.
// This data is retrieved directly from the in-memory `RolesIndex` in the `Context`,
// providing quick access to role metadata without reading full role definitions from disk.
func (w *Workspace) ListRoles() ([]RoleSummary, error) {
	roles := make([]RoleSummary, 0, len(w.Context.Indexes.RolesIndex))
	for _, r := range w.Context.Indexes.RolesIndex {
		roles = append(roles, r)
	}
	return roles, nil
}

// ListPreferences returns a slice of all preference summaries.
// This data is retrieved directly from the in-memory `PreferencesIndex` in the `Context`,
// enabling efficient listing of user preferences.
func (w *Workspace) ListPreferences() ([]PreferenceSummary, error) {
	preferences := make([]PreferenceSummary, 0, len(w.Context.Indexes.PreferencesIndex))
	for _, p := range w.Context.Indexes.PreferencesIndex {
		preferences = append(preferences, p)
	}
	return preferences, nil
}


// loadRole loads a role by its name from `roles/<name>.json`.
// This is an internal helper function.
func (w *Workspace) loadRole(name string) (Role, error) {
	rolePath := filepath.Join(w.RootDir, "roles", fmt.Sprintf("%s.json", name))
	var role Role
	data, err := os.ReadFile(rolePath)
	if err != nil {
		return Role{}, fmt.Errorf("failed to read role file %s: %w", name, err)
	}
	if err := json.Unmarshal(data, &role); err != nil {
		return Role{}, fmt.Errorf("failed to parse role data from %s: %w", name, err)
	}
	return role, nil
}

// SavePreference saves a user preference to `preferences/<id>.json`.
// After saving the file, it updates the `PreferencesIndex` in the `Context`
// and persists the updated `Context` to disk.
func (w *Workspace) SavePreference(pref Preference) error {
	prefPath := filepath.Join(w.RootDir, "preferences", fmt.Sprintf("%s.json", pref.ID))
	if err := w.writeJSON(prefPath, pref); err != nil {
		return fmt.Errorf("failed to save preference %s: %w", pref.ID, err)
	}

	snippet := pref.Content
	if len(snippet) > 100 { // Limit snippet length for display in summary
		snippet = snippet[:100] + "..."
	}
	w.Context.Indexes.PreferencesIndex[pref.ID] = PreferenceSummary{
		ID:             pref.ID,
		Timestamp:      pref.Timestamp,
		ContentSnippet: snippet,
	}
	if err := w.saveContext(w.Context); err != nil {
		return fmt.Errorf("failed to update context after saving preference: %w", err)
	}
	return w.logAction(fmt.Sprintf("Saved preference %s", pref.ID))
}

// LoadPreference loads a single preference by its unique ID from `preferences/<id>.json`.
// It returns a pointer to the `Preference` struct or an error if the file
// cannot be read or parsed.
func (w *Workspace) LoadPreference(id string) (*Preference, error) {
	prefPath := filepath.Join(w.RootDir, "preferences", fmt.Sprintf("%s.json", id))
	data, err := os.ReadFile(prefPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read preference %s: %w", id, err)
	}
	var pref Preference
	if err := json.Unmarshal(data, &pref); err != nil {
		return nil, fmt.Errorf("failed to parse preference %s: %w", id, err)
	}
	return &pref, nil
}

// DeletePreference deletes a preference file from `preferences/<id>.json`
// and removes its entry from the `PreferencesIndex` in the `Context`.
// The updated `Context` is then saved to disk.
func (w *Workspace) DeletePreference(id string) error {
	prefPath := filepath.Join(w.RootDir, "preferences", fmt.Sprintf("%s.json", id))
	if err := os.Remove(prefPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete preference file %s: %w", id, err)
	}

	delete(w.Context.Indexes.PreferencesIndex, id)
	if err := w.saveContext(w.Context); err != nil {
		return fmt.Errorf("failed to update context after deleting preference: %w", err)
	}
	return w.logAction(fmt.Sprintf("Deleted preference %s", id))
}


// loadSession loads the current active session from `session.json`.
// This internal helper function is responsible for reading the session file,
// unmarshaling its data, and then loading the full `Role` configuration
// (which is only stored by name in the session JSON) to ensure the in-memory
// `Session` struct is complete.
func (w *Workspace) loadSession() (*Session, error) {
	sessionPath := filepath.Join(w.RootDir, "session.json")
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no active session found at %s: %w", sessionPath, err)
	}

	data, err := os.ReadFile(sessionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read active session file %s: %w", sessionPath, err)
	}

	var session Session
	// Note: Session.UnmarshalJSON will only populate the Role.Name initially
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse active session data from %s: %w", sessionPath, err)
	}

	// Now, load the full role data using the name unmarshaled from session.json.
	// This step is crucial because Session.UnmarshalJSON only gets the role's name.
	role, err := w.loadRole(session.Role.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to load full role data for session %s (role name: %s): %w", session.ID, session.Role.Name, err)
	}
	session.Role = role // Assign the fully loaded role to the session

	return &session, nil
}

// loadContext loads the `context.json` file into the Workspace's `Context` field.
// This is an internal helper function.
func (w *Workspace) loadContext() error {
	contextPath := filepath.Join(w.RootDir, "context.json")
	data, err := os.ReadFile(contextPath)
	if err != nil {
		return fmt.Errorf("failed to read context: %w", err)
	}
	var context Context
	if err := json.Unmarshal(data, &context); err != nil {
		return fmt.Errorf("failed to parse context: %w", err)
	}
	w.Context = context
	return nil
}

// saveContext saves the current Workspace's `Context` to `context.json`.
// This is an internal helper function, typically called after any modifications
// to the `Context` (including its indexes) to persist changes.
func (w *Workspace) saveContext(context Context) error {
	return w.writeJSON(filepath.Join(w.RootDir, "context.json"), context)
}

// saveRole saves an AI role configuration to `roles/<name>.json`.
// After saving the role file, it updates the `RolesIndex` in the `Context`
// and persists the updated `Context` to disk.
func (w *Workspace) saveRole(role Role) error {
	rolePath := filepath.Join(w.RootDir, "roles", fmt.Sprintf("%s.json", role.Name))
	if err := w.writeJSON(rolePath, role); err != nil {
		return fmt.Errorf("failed to save role %s: %w", role.Name, err)
	}

	// Update the index
	w.Context.Indexes.RolesIndex[role.Name] = RoleSummary{
		Name:        role.Name,
		Label:       role.Label,
		Description: role.Description,
	}
	if err := w.saveContext(w.Context); err != nil {
		return fmt.Errorf("failed to update context after saving role: %w", err)
	}
	return w.logAction(fmt.Sprintf("Saved role %s", role.Name))
}

// DeleteRole deletes a role file from `roles/<name>.json` and removes its entry
// from the `RolesIndex` in the `Context`. The updated `Context` is then saved to disk.
func (w *Workspace) DeleteRole(name string) error {
	rolePath := filepath.Join(w.RootDir, "roles", fmt.Sprintf("%s.json", name))
	if err := os.Remove(rolePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete role file %s: %w", name, err)
	}

	delete(w.Context.Indexes.RolesIndex, name)
	if err := w.saveContext(w.Context); err != nil {
		return fmt.Errorf("failed to update context after deleting role: %w", err)
	}
	return w.logAction(fmt.Sprintf("Deleted role %s", name))
}


// saveSession saves the given `Session` struct to the active `session.json` file.
// This is an internal helper function.
func (w *Workspace) saveSession(session Session) error {
	return w.writeJSON(filepath.Join(w.RootDir, "session.json"), session)
}


// writeJSON is a utility helper function that writes data to a JSON file.
// It ensures proper indentation (2 spaces) and file permissions (0644 - owner rw, group r, others r).
// This is an internal helper function used by various save operations.
func (w *Workspace) writeJSON(path string, data interface{}) error {
	// 0644: owner rw, group r, others r
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close() // Ensure the file is closed even if an error occurs

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Use 2 spaces for indentation
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write JSON to %s: %w", path, err)
	}
	return nil
}

// logAction appends a timestamped action entry to the daily log file.
// Log files are stored in the `logs/` subdirectory, named by date (e.g., `2006-01-02.log`).
// This is an internal helper function for logging operational events within the workspace.
func (w *Workspace) logAction(action string) error {
	logDir := filepath.Join(w.RootDir, "logs")
	logFile := filepath.Join(logDir, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))) // e.g., 2024-07-30.log

	// Open file in append mode, create if it doesn't exist, write-only
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	logEntry := fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), action)
	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}
	return nil
}

