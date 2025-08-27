// Package tmux provides utilities for interacting with and managing tmux sessions.
package tmux

import (
	"fmt"

	"github.com/phanorcoll/muxie/internal/config"
)

// StartSession creates a new tmux session with the given sessionName and starting directory.
// It then creates the specified windows and panes, running the configured commands in each pane.
// Returns an error if any tmux operation fails.
func StartSession(sessionName, dir string, windows []config.Window) error {
	if err := CreateSession(sessionName, dir); err != nil {
		return fmt.Errorf("failed to create new session '%s': %w", sessionName, err)
	}
	directory := expandHomeDir(dir)
	for _, w := range windows {
		if err := NewWindow(sessionName, w.Name, directory); err != nil {
			return fmt.Errorf("failed to create window '%s' in session '%s': %w", w.Name, sessionName, err)
		}

		for j, p := range w.Panes {
			if j > 0 {
				if err := SplitWindow(sessionName, w.Name, w.Layout); err != nil {
					return fmt.Errorf("failed to split window '%s' in session '%s': %w", w.Name, sessionName, err)
				}
			}
			if err := SendKeys(sessionName, w.Name, j, fmt.Sprintf("cd %s && clear && %s", directory, p.Command)); err != nil {
				return fmt.Errorf("failed to send keys to pane %d in window '%s' of session '%s': %w", j, w.Name, sessionName, err)
			}
		}
	}

	return nil
}
