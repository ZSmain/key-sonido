import json
import logging
import sys
from threading import Thread
from pynput import keyboard
import pygame

# Load configuration
with open("config.json", "r") as f:
    config = json.load(f)

# Initialize logging
logging.basicConfig(filename="keyboard_sound.log", level=logging.ERROR)

# Initialize pygame mixer
pygame.mixer.init()
pygame.mixer.music.set_volume(config.get("volume", 0.5))

# Load sounds
sounds = {
    "alphanumeric": pygame.mixer.Sound(config["sounds"]["alphanumeric"]),
    "space": pygame.mixer.Sound(config["sounds"]["space"]),
    "enter": pygame.mixer.Sound(config["sounds"]["enter"]),
    "backspace": pygame.mixer.Sound(config["sounds"]["backspace"]),
}


def play_sound(key):
    try:
        if key == keyboard.Key.space:
            sounds["space"].play()
        elif key == keyboard.Key.enter:
            sounds["enter"].play()
        elif key == keyboard.Key.backspace:
            sounds["backspace"].play()
        elif hasattr(key, "char") and key.char.isalnum():
            sounds["alphanumeric"].play()
    except Exception as e:
        logging.error(f"Error playing sound for key {key}: {e}")


def on_press(key):
    play_sound(key)


def start_listener():
    with keyboard.Listener(on_press=on_press) as listener:
        listener.join()


def main():
    if len(sys.argv) != 2:
        print("Usage: ./main.py start|stop")
        sys.exit(1)

    command = sys.argv[1]
    if command == "start":
        thread = Thread(target=start_listener)
        thread.daemon = True
        thread.start()
        print("Keyboard sound effects started.")
        thread.join()
    elif command == "stop":
        print("Keyboard sound effects stopped.")
        sys.exit()
    else:
        print("Unknown command. Use start or stop.")
        sys.exit(1)


if __name__ == "__main__":
    main()
