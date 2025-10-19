// Package tui contains terminal user interface commands and related functionality.
package tui

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/phanorcoll/muxie/internal/config"
	"github.com/phanorcoll/muxie/internal/tmux"
)

// sessionsResponseMsg is a message containing the list of tmux sessions,
// the currently active session, and any error encountered when retrieving them.
type sessionsResponseMsg struct {
	SessionsList  []list.Item // List of tmux sessions
	ActiveSession string
	Err           error
}

// activeResponseMsg represents a message containing the name of the currently active session.
type activeResponseMsg struct {
	ActiveSession string
}

// getSessionsCmd retrieves the list of tmux sessions and the currently active session.
// It returns a sessionsResponseMsg containing the sessions list, the active session,
// and any error encountered during retrieval.
func getSessionsCmd(config *config.Config) tea.Cmd {
	sl, err := tmux.GetSessionsList()
	if err != nil {
		return func() tea.Msg {
			return sessionsResponseMsg{
				Err: err,
			}
		}
	}
	var sessions []list.Item
	activeSession, err := tmux.GetActiveSession()

	for _, sessionInfo := range config.Sessions {
		sessions = append(sessions, session{
			sessionName:   sessionInfo.Name,
			numWindows:    len(sessionInfo.Windows),
			activeSession: activeSession,
			isFromConfig:  true,
			isRunning:     false,
		})
	}

	existingSessions := make(map[string]bool)
	for i, s := range sessions {
		for _, runningSession := range sl {
			if s.(session).sessionName == runningSession.Name {
				sessionItem := s.(session)
				sessionItem.isRunning = true
				sessionItem.numWindows = runningSession.NumberWindows
				sessions[i] = sessionItem
				existingSessions[sessionItem.sessionName] = true
			}
		}
	}

	for _, sessionInfo := range sl {
		if !existingSessions[sessionInfo.Name] {
			sessions = append(sessions, session{
				sessionName:   sessionInfo.Name,
				numWindows:    sessionInfo.NumberWindows,
				activeSession: activeSession,
				isFromConfig:  false,
				isRunning:     true,
			})
		}
	}

	if err != nil {
		log.Println("Error getting active session:", err)
		return func() tea.Msg {
			return sessionsResponseMsg{
				Err: err,
			}
		}
	}
	return func() tea.Msg {
		return sessionsResponseMsg{
			SessionsList:  moveActiveSessionToTop(sessions, activeSession),
			ActiveSession: activeSession,
		}
	}
}

// moveActiveSessionToTop reorders the sessions so that the active session is first,
// followed by other running sessions, then the rest.
func moveActiveSessionToTop(sessions []list.Item, activeSession string) []list.Item {
	var active []list.Item
	var running []list.Item
	var others []list.Item

	for _, s := range sessions {
		sess, ok := s.(session)
		if !ok {
			others = append(others, s)
			continue
		}
		if sess.sessionName == activeSession {
			active = append(active, s)
		} else if sess.isRunning {
			running = append(running, s)
		} else {
			others = append(others, s)
		}
	}


	combined := append(append(active, running...), others...)

	// Find the session that shows up right before the "others" and
	// set it's addSpacingUnder property to true.
	// This is then interpreted by the Render function to add a single
	// line of white space underneath. This helps separate the active/running
	// sessions from the unstarted config sessions.
	if 0 < len(others) && len(others) < len(combined) {
		sess := combined[len(combined) - len(others) - 1].(session)
		sess.addSpacingUnder = true
		combined[len(combined) - len(others) - 1] = sess
	}

	return combined
}
