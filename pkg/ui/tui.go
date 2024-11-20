package ui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zsmain/key-sonido/config"
)

type state int

const (
	stateMenu state = iota
	stateSelectProfile
	stateAdjustVolume
)

type model struct {
	state         state
	menuItems     []string
	selectedIndex int
	config        *config.Config

	// Fields for profile selection
	profileList  []string
	profileIndex int

	// Fields for volume adjustment
	volumeInput textinput.Model
}

func NewModel(cfg *config.Config) model {
	return model{
		state:     stateMenu,
		menuItems: []string{"Sound Profile", "Volume", "Toggle Sound", "Exit"},
		config:    cfg,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateMenu:
		// Handle menu state
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.selectedIndex > 0 {
					m.selectedIndex--
				}
			case "down", "j":
				if m.selectedIndex < len(m.menuItems)-1 {
					m.selectedIndex++
				}
			case "enter":
				switch m.menuItems[m.selectedIndex] {
				case "Sound Profile":
					// Transition to select profile state
					m.state = stateSelectProfile
					// Load profiles
					profiles, err := os.ReadDir("assets")
					if err != nil {
						// Handle error
						return m, nil
					}
					for _, profile := range profiles {
						if profile.IsDir() {
							m.profileList = append(m.profileList, profile.Name())
						}
					}
				case "Volume":
					// Transition to adjust volume state
					m.state = stateAdjustVolume
					m.volumeInput = textinput.New()
					m.volumeInput.Placeholder = "Enter volume (0-100)"
					m.volumeInput.SetValue(fmt.Sprintf("%d", m.config.Volume))
					m.volumeInput.CharLimit = 3
					m.volumeInput.Focus()
				case "Toggle Sound":
					// Toggle sound
					m.config.Enabled = !m.config.Enabled
					config.SaveConfig(m.config)
				case "Exit":
					return m, tea.Quit
				}
			case "q", "esc":
				return m, tea.Quit
			}
		}
	case stateSelectProfile:
		// Handle profile selection state
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.profileIndex > 0 {
					m.profileIndex--
				}
			case "down", "j":
				if m.profileIndex < len(m.profileList)-1 {
					m.profileIndex++
				}
			case "enter":
				// Set selected profile
				m.config.SoundProfile = m.profileList[m.profileIndex]
				config.SaveConfig(m.config)
				// Return to menu
				m.state = stateMenu
			case "q", "esc":
				// Return to menu without changing
				m.state = stateMenu
			}
		}
	case stateAdjustVolume:
		// Handle volume adjustment state
		var cmd tea.Cmd
		m.volumeInput, cmd = m.volumeInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Save volume
				val := m.volumeInput.Value()
				vol, err := strconv.Atoi(val)
				if err == nil && vol >= 0 && vol <= 100 {
					m.config.Volume = vol
					config.SaveConfig(m.config)
				}
				// Return to menu
				m.state = stateMenu
			case "q", "esc":
				// Return to menu without changing
				m.state = stateMenu
			}
		}
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateMenu:
		// Render the main menu
		s := TitleStyle.Render("KeySonido Configuration") + "\n\n"
		for i, item := range m.menuItems {
			if i == m.selectedIndex {
				s += SelectedMenuItemStyle.Render(item) + "\n"
			} else {
				s += MenuItemStyle.Render(item) + "\n"
			}
		}
		return s
	case stateSelectProfile:
		// Render the profile selection menu
		s := TitleStyle.Render("Select Sound Profile") + "\n\n"
		for i, profile := range m.profileList {
			if i == m.profileIndex {
				s += SelectedMenuItemStyle.Render(profile) + "\n"
			} else {
				s += MenuItemStyle.Render(profile) + "\n"
			}
		}
		return s
	case stateAdjustVolume:
		// Render the volume adjustment input
		s := TitleStyle.Render("Adjust Volume") + "\n\n"
		s += m.volumeInput.View() + "\n"
		s += "\nPress Enter to save or Esc to cancel."
		return s
	default:
		return "Unknown state"
	}
}

func StartTUI(cfg *config.Config) {
	p := tea.NewProgram(NewModel(cfg))
	if err := p.Start(); err != nil {
		fmt.Println("Error running TUI:", err)
	}
}
