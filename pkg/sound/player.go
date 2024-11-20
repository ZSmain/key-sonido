package sound

import (
	"os"
	"time"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func playSound(filePath string, volume int) {
	// Open the sound file
	f, err := os.Open(filePath)
	if err != nil {
		// Handle error
		return
	}
	defer f.Close()
	// Decode the .wav file
	streamer, format, err := wav.Decode(f)
	if err != nil {
		// Handle error
		return
	}
	defer streamer.Close()
	// Initialize the speaker with the sample rate
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// Adjust volume if needed
	volumeControl := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   float64(volume-100) / 100, // Convert volume to acceptable range
		Silent:   false,
	}
	// Play the sound
	speaker.Play(volumeControl)
}
