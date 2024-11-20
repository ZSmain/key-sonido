package sound

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/zsmain/key-sonido/config"
)

func PlayRandomSound(cfg *config.Config) {
	// Get the path to the current sound profile
	profilePath := filepath.Join("assets", cfg.SoundProfile)

	// List all subdirectories in the profile path
	dirs, err := os.ReadDir(profilePath)
	if err != nil {
		// Handle error
		return
	}

	// Collect the names of subdirectories
	var subfolders []string
	for _, dir := range dirs {
		if dir.IsDir() {
			subfolders = append(subfolders, dir.Name())
		}
	}

	// If no subfolders are found, proceed with the profile path
	var searchPaths []string
	if len(subfolders) > 0 {
		// Create a local random number generator
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		// Select a random subfolder
		randomSubfolder := subfolders[rng.Intn(len(subfolders))]
		// Path to the random subfolder
		searchPaths = append(searchPaths, filepath.Join(profilePath, randomSubfolder, "*.wav"))
	} else {
		// No subfolders found, use the profile path
		searchPaths = append(searchPaths, filepath.Join(profilePath, "*.wav"))
	}

	// Get list of sound files from the search paths
	var soundFiles []string
	for _, pattern := range searchPaths {
		files, err := filepath.Glob(pattern)
		if err == nil {
			soundFiles = append(soundFiles, files...)
		}
	}

	if len(soundFiles) == 0 {
		// Handle no sound files found
		return
	}

	// Select a random sound file using the global rand
	randomFile := soundFiles[rand.Intn(len(soundFiles))]

	// Call playSound with the selected file
	playSound(randomFile, cfg.Volume)
}
