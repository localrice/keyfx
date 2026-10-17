package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/holoplot/go-evdev"
)

func main() {
	const soundFile = "gun_shot.mp3"

	keys, err := keyboardEvents()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(soundFile)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	fmt.Println("Buffer length:", buffer.Len())
	streamer.Close()
	f.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// speaker.PlayAndWait(streamer)
	// speaker.PlayAndWait(buffer.Streamer(0, buffer.Len()))

	for range keys {
		fmt.Println("key press")
		// speaker.Play(beep.Clone(streamer))
		speaker.Play(buffer.Streamer(0, buffer.Len()))
	}

}

func keyboardEvents() (<-chan struct{}, error) {
	devices, err := os.ReadDir("/dev/input")
	if err != nil {
		return nil, err
	}

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

		go func() {
			defer dev.Close()

			for {
				event, err := dev.ReadOne()
				if err != nil {
					return
				}

				if event.Type == evdev.EV_KEY && event.Value == 1 {
					keys <- struct{}{}
				}
			}
		}()
	}

	if !found {
		close(keys)
		return nil, fmt.Errorf("no accessible input devices found")
	}

	return keys, nil
}
