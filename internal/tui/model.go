// Package tui contains terminal user interface commands and related functionality.
package tui

import (
	"github.com/phanorcoll/muxie/internal/config"
	"github.com/phanorcoll/muxie/internal/log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// statusData holds information about a status indicator in the TUI, including its icon,
// action, color, and a human-readable action title.
type statusData struct {
	icon        string
	action      string
	color       string
	actionTitle string
}

// Model represents the main state of the TUI application.
// It contains configuration, session data, UI state, version, and input models.
type Model struct {
	version       string          // Application version
	config        *config.Config  // Application configuration
	statusData    statusData      // Current status indicator data
	activeSession string          // Name of the currently active session
	showInput     bool            // Whether the input field is visible
	logger        log.Logger      // Logger for debugging
	keys          keyMap          // Key bindings for the TUI
	sessionList   list.Model      // List model for displaying sessions
	help          help.Model      // Help model for displaying key bindings/help
	sessionInput  textinput.Model // Text input model for session creation/renaming
}

// NewModel creates and returns a new Model instance for the TUI application.
// It initializes the session list, status data, key bindings, help model, and session input field.
func NewModel(cfg *config.Config, logger log.Logger, version string) Model {
	newSessionInput := textinput.New()
	newSessionInput.CharLimit = 50
	newSessionInput.Width = 20
	sl := initList()
	return Model{
		version: version,
		config:  cfg,
		logger:  logger,
		statusData: statusData{
			icon:   "",
			action: "",
			color:  "#FFF7DB",
		},
		sessionList:  sl,
		keys:         defaultKeyMap,
		help:         help.New(),
		sessionInput: newSessionInput,
	}
}

// Init is part of the Bubble Tea Model interface and initializes the program.
// It returns an initial command to run, or nil if there is none.
func (m Model) Init() tea.Cmd {
	return getSessionsCmd(m.config)
}
