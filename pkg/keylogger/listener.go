package keylogger

import (
	"fmt"

	"github.com/MarinX/keylogger"
	"github.com/zsmain/key-sonido/config"
	"github.com/zsmain/key-sonido/pkg/sound"
)

func StartKeyListener(cfg *config.Config) {
	// Find keyboard device
	keyboard := keylogger.FindKeyboardDevice()
	if keyboard == "" {
		fmt.Println("No keyboard found...")
		return
	}
	// Initialize keylogger
	kl, err := keylogger.New(keyboard)
	if err != nil {
		fmt.Println("Error initializing keylogger:", err)
		return
	}
	// Read key events
	events := kl.Read()
	for e := range events {
		if e.Type == keylogger.EvKey && e.KeyPress() {
			if isToggleHotkey(e) {
				// Toggle sound on/off
				cfg.Enabled = !cfg.Enabled
				config.SaveConfig(cfg)
				continue
			}
			if cfg.Enabled {
				// On keypress, call PlayRandomSound from the sound package
				sound.PlayRandomSound(cfg)
			}
		}
	}
}

func isToggleHotkey(e keylogger.InputEvent) bool {
	// Implement hotkey detection, e.g., Ctrl + Shift + K
	// ...code to determine if the event matches the hotkey...
	return false
}
