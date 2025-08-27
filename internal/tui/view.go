// Package tui contains terminal user interface commands and related functionality.
package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	widthTerm, heightTerm = func() (int, int) {
		w, h, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return 80, 20
		}
		return w, h
	}()

	width  = 47
	height = 13
	// set of colors
	normal  = lipgloss.Color("#EEEEEE")
	subtle  = lipgloss.Color("#72726F")
	special = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	// set of styles
	inputStyle    = lipgloss.NewStyle()
	base          = lipgloss.NewStyle().Foreground(normal)
	url           = lipgloss.NewStyle().Foreground(special).Render
	versionStyle  = lipgloss.NewStyle().Foreground(subtle).Render
	titleStyle    = lipgloss.NewStyle().Italic(true).Bold(true).Foreground(normal)
	bulletDivider = lipgloss.NewStyle().
			SetString("•").
			Padding(0, 1).
			Foreground(subtle).
			String()
	divider = base.
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(subtle).
		Width(width).
		Render()
	borderStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#5a5a5a"))
	inputHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(subtle))

	dialogBoxStyle = lipgloss.NewStyle().
			Padding(1, 2, 1).
			Border(lipgloss.NormalBorder())
)

// View renders the entire TUI based on the current state of the Model.
// It constructs the header, content area (input dialog or session list), and applies styling.
// The output is centered and bordered according to the terminal dimensions.
func (m Model) View() string {
	doc := strings.Builder{}
	w := lipgloss.Width
	muxiemeta := versionStyle(m.version)
	// header
	{
		// versionInfo := versionStyle(m.version)
		logo := titleStyle.Render("󱌖 Muxie")
		logoversion := lipgloss.JoinHorizontal(lipgloss.Right, logo, " ", muxiemeta)
		actionIcon := base.Foreground(lipgloss.Color(m.statusData.color)).Render(m.statusData.icon)
		activeSesh := titleStyle.Width(width - w(logoversion)).Render(actionIcon + bulletDivider + m.activeSession)
		row := lipgloss.JoinHorizontal(lipgloss.Center, activeSesh, logoversion)
		header := lipgloss.JoinVertical(lipgloss.Top, row, divider)
		doc.WriteString(header)
	}
	// content
	{
		if m.showInput {
			question := base.Foreground(lipgloss.Color(m.statusData.color)).Render(m.statusData.actionTitle)
			input := inputStyle.Foreground(lipgloss.Color(m.statusData.color)).Render(m.sessionInput.View())
			inputHelpStyle := inputHelpStyle.Render("esc - cancel")
			ui := lipgloss.JoinVertical(lipgloss.Left, question, input, inputHelpStyle)

			dialog := lipgloss.Place(width, 9,
				lipgloss.Center, lipgloss.Center,
				dialogBoxStyle.BorderForeground(lipgloss.Color(m.statusData.color)).Render(ui),
			)
			doc.WriteString(dialog)
		} else {
			doc.WriteString(m.sessionList.View())
			doc.WriteString("\n" + m.help.View(m.keys))
		}
	}
	content := doc.String()
	borderedContainer := borderStyle.Width(width).Height(height).Render(content)
	return lipgloss.Place(
		widthTerm, heightTerm, lipgloss.Center, lipgloss.Center, borderedContainer,
	)
}
