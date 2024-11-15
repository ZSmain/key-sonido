use device_query::{DeviceQuery, DeviceState, Keycode};
use inquire::validator::Validation;
use inquire::{CustomType, Select};
use rodio::{Decoder, OutputStream, Sink};
use std::fs::File;
use std::io::BufReader;
use std::sync::{Arc, Mutex};
use std::{thread, time::Duration};

struct Config {
    sound_file: String,
    volume: f32,
    is_running: bool,
}

impl Config {
    fn new() -> Self {
        Self {
            sound_file: "sounds/click1.wav".to_string(),
            volume: 0.5,
            is_running: false,
        }
    }
}

fn main() {
    let config = Arc::new(Mutex::new(Config::new()));

    loop {
        let options = vec![
            "Start Service",
            "Stop Service",
            "Select Sound",
            "Adjust Volume",
            "Exit",
        ];

        let choice = Select::new("Choose an option:", options).prompt().unwrap();

        match choice {
            "Start Service" => {
                let cfg = config.clone();
                {
                    let mut config_lock = cfg.lock().unwrap();
                    if !config_lock.is_running {
                        config_lock.is_running = true;
                        println!("Service started.");
                        let cfg = cfg.clone();
                        thread::spawn(move || {
                            run_service(cfg);
                        });
                    } else {
                        println!("Service is already running.");
                    }
                }
            }
            "Stop Service" => {
                let mut config_lock = config.lock().unwrap();
                if config_lock.is_running {
                    config_lock.is_running = false;
                    println!("Service stopping...");
                } else {
                    println!("Service is not running.");
                }
            }
            "Select Sound" => {
                let sounds = vec!["click1.wav", "click2.wav", "click3.wav"];
                let sound = Select::new("Select a sound:", sounds).prompt().unwrap();
                let mut cfg = config.lock().unwrap();
                cfg.sound_file = format!("sounds/{}", sound);
                println!("Sound selected: {}", cfg.sound_file);
            }
            "Adjust Volume" => {
                let volume: f32 = CustomType::new("Enter volume (0.0 - 1.0):")
                    .with_error_message("Please enter a number between 0.0 and 1.0")
                    .with_validator(|input: &f32| {
                        if *input >= 0.0 && *input <= 1.0 {
                            Ok(Validation::Valid)
                        } else {
                            Ok(Validation::Invalid(
                                "Volume must be between 0.0 and 1.0".into(),
                            ))
                        }
                    })
                    .prompt()
                    .unwrap();

                let mut cfg = config.lock().unwrap();
                cfg.volume = volume;
                println!("Volume set to: {}", volume);
            }
            "Exit" => {
                let mut config_lock = config.lock().unwrap();
                config_lock.is_running = false;
                println!("Exiting.");
                break;
            }
            _ => (),
        }
    }
}

fn run_service(config: Arc<Mutex<Config>>) {
    let device_state = DeviceState::new();
    let mut last_keys = device_state.get_keys();

    let (_stream, stream_handle) = OutputStream::try_default().unwrap();

    while config.lock().unwrap().is_running {
        let keys = device_state.get_keys();

        if keys != last_keys {
            let new_keys: Vec<&Keycode> = keys.iter().filter(|k| !last_keys.contains(k)).collect();
            if !new_keys.is_empty() {
                let cfg = config.lock().unwrap();
                play_sound(&stream_handle, &cfg.sound_file, cfg.volume);
            }
            last_keys = keys;
        }

        thread::sleep(Duration::from_millis(10));
    }

    println!("Service stopped.");
}

fn play_sound(stream_handle: &rodio::OutputStreamHandle, file_path: &str, volume: f32) {
    if let Ok(file) = File::open(file_path) {
        let source = Decoder::new(BufReader::new(file)).unwrap();
        let sink = Sink::try_new(stream_handle).unwrap();
        sink.set_volume(volume);
        sink.append(source);
        sink.detach();
    } else {
        eprintln!("Error: Could not play sound {}", file_path);
    }
}
