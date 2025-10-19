// Package tui contains terminal user interface commands and related functionality.
package tui

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listHeight   = 10
	defaultWidth = 40
)

var (
	itemStyle              = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#7C7C7E"))
	selectedItemStyle      = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#7fd1ae")).Bold(true)
	paginationStyle        = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	statusMessageStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).Render
	activeSessionStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F2CC4A"))
	activeSessionHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#444341")).Render
	errorStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render
)

// session represents a tmux session in the TUI.
// sessionName: the name of the session.
// numWindows: the number of windows in the session.
// activeSession: the currently active session, used for highlighting.
// isFromConfig: true if the session is defined in the config file.
// isRunning: true if the session is currently running.
type session struct {
	sessionName     string
	numWindows      int
	activeSession   string
	isFromConfig    bool
	isRunning       bool
	addSpacingUnder bool
}

func (i session) SessionName() string { return i.sessionName }
func (i session) FilterValue() string { return i.sessionName }

type sessionDelegate struct{}

func (d sessionDelegate) Height() int                             { return 1 }
func (d sessionDelegate) Spacing() int                            { return 0 }
func (d sessionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d sessionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(session)
	if !ok {
		return
	}

	var icon string
	if i.isFromConfig {
		icon = "󰗼 "
	} else {
		icon = "󰆍 "
	}

	name := formatSessionName(i.sessionName)
	activesession := formatSessionName(i.activeSession)

	if !i.isRunning {
		icon = "󰄜 "
		name = activeSessionHelpStyle(name)
	}

	if name == activesession {
		name = activeSessionStyle.Render(name + " " + activeSessionHelpStyle("  󰞓 active"))
	}

	desc := fmt.Sprintf("-%d󱂬 - ", i.numWindows) + icon + name

	// This just adds some separation between active/running sessions
	// and the config sessions that have not been started yet
	if i.addSpacingUnder {
		desc += "\n"
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(" " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(desc))
}

func formatSessionName(name string) string {
	if len(name) == 0 {
		return ""
	}
	name = strings.ToLower(name)
	r := []rune(name)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// initList initializes and returns a new list.Model with default settings
// for use in the TUI. The list is configured to hide the title, status bar,
// and help, enables filtering, and disables quit keybindings.
func initList() list.Model {
	sessionList := list.New([]list.Item{}, sessionDelegate{}, defaultWidth, listHeight)
	sessionList.Styles.Title = titleStyle
	sessionList.Title = "Sessions"
	sessionList.SetShowTitle(true)
	sessionList.SetShowStatusBar(false)
	sessionList.SetShowHelp(false)
	sessionList.SetFilteringEnabled(true)
	sessionList.DisableQuitKeybindings()
	sessionList.StatusMessageLifetime = time.Second * 2
	return sessionList
}
