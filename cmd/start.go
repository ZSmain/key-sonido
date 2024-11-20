package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zsmain/key-sonido/config"
	"github.com/zsmain/key-sonido/pkg/keylogger"
	"github.com/zsmain/key-sonido/pkg/ui"
)

var daemonMode bool

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start KeySonido",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if config directory exists, if not create it
		if _, err := os.Stat("config"); os.IsNotExist(err) {
			err = os.Mkdir("config", 0755)
			if err != nil {
				fmt.Println("Error creating config directory:", err)
				return
			}
		}

		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			if os.IsNotExist(err) {
				// Create default configuration if it doesn't exist
				cfg = &config.Config{
					SoundProfile: "click1",
					Volume:       80,
				}
				err = config.SaveConfig(cfg)
				if err != nil {
					fmt.Println("Error saving default configuration:", err)
					return
				}
			} else {
				fmt.Println("Error loading configuration:", err)
				return
			}
		}

		if daemonMode {
			// Run in background
			// ...code to daemonize the process...
		}

		// Start key listener
		go keylogger.StartKeyListener(cfg)

		// Start interactive CLI
		ui.StartTUI(cfg)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&daemonMode, "daemon", "d", false, "Run as background process")
}
