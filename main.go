package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/holoplot/go-evdev"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Sound string `yaml:"sound"`
}

func main() {
	if err := runDaemon(); err != nil {
		log.Fatal(err)
	}
}

func runDaemon() error {
	const soundFile = "gun_shot.mp3"

	keys, err := keyboardEvents()
	if err != nil {
		return err
	}

	f, err := os.Open(soundFile)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return err
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	streamer.Close()
	f.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// speaker.PlayAndWait(streamer)
	// speaker.PlayAndWait(buffer.Streamer(0, buffer.Len()))

	// play sound from the keys channel provided by keyboardEvents
	for range keys {
		speaker.Play(buffer.Streamer(0, buffer.Len()))
	}
	return nil
}

func keyboardEvents() (<-chan struct{}, error) {
	devices, err := os.ReadDir("/dev/input")
	if err != nil {
		return nil, err
	}

	// send an empty struct as a signal that a key was pressed
	// in the future, maybe send the key information to signify which key was pressed
	keys := make(chan struct{})

	var found bool

	for _, device := range devices {
		//example: /dev/input/event*
		if !strings.HasPrefix(device.Name(), "event") {
			continue
		}

		path := "/dev/input/" + device.Name()

		dev, err := evdev.Open(path)
		if err != nil {
			continue
		}

		found = true

		go func(dev *evdev.InputDevice) {
			defer dev.Close()

			for {
				event, err := dev.ReadOne()
				if err != nil {
					return
				}

				//EV_KEY with values:
				//	0 - key released
				//	1 - key pressed
				//	2 - key held
				if event.Type == evdev.EV_KEY && event.Value == 1 {
					keys <- struct{}{}
				}
			}
		}(dev)
	}

	if !found {
		close(keys)
		return nil, fmt.Errorf("no accessible input devices found")
	}

	return keys, nil
}

func loadConfig() (Config, error) {
	var cfg Config

	configDir, err := os.UserConfigDir()
	if err != nil {
		return cfg, err
	}

	configPath := filepath.Join(configDir, "keyfx", "config.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
