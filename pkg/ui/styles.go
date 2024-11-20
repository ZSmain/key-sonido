package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Define styles for TUI components

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4"))

	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("#FAFAFA"))

	SelectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#7D56F4")).
				Background(lipgloss.Color("#FAFAFA"))

	// ...additional styles...
)
