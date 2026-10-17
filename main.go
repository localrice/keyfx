package main

import (
	"fmt"
	"log"
	"sync"
	"syscall"

	evdev "github.com/holoplot/go-evdev"
)

func main() {
	// f, err := os.Open("gun_shot.mp3")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// streamer, format, err := mp3.Decode(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer streamer.Close()
	// sr := format.SampleRate / 100
	// speaker.Init(sr, format.SampleRate.N(time.Second/10))
	// speaker.PlayAndWait(streamer)

	// select {}

	devices, err := evdev.ListDevicePaths()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	found := 0

	for _, path := range devices {
		dev, err := evdev.OpenWithFlags(path.Path, syscall.O_RDONLY)
		if err != nil {
			continue
		}

		// Identify keyboards by checking for ordinary key codes.
		keys := dev.CapableEvents(evdev.EV_KEY)
		keyboard := false

		for _, code := range keys {
			if code == evdev.KEY_A || code == evdev.KEY_SPACE {
				keyboard = true
				break
			}
		}

		if !keyboard {
			dev.Close()
			continue
		}

		name, _ := dev.Name()
		fmt.Println("Listening:", name, "-", path.Path)
		found++

		wg.Add(1)
		go func(d *evdev.InputDevice) {
			defer wg.Done()
			defer d.Close()

			for {
				event, err := d.ReadOne()
				if err != nil {
					return
				}

				// EV_KEY + value 1 = key-down.
				if event.Type == evdev.EV_KEY && event.Value == 1 {
					fmt.Printf("%s\n", event.CodeName())
				}
			}
		}(dev)
	}

	if found == 0 {
		log.Fatal("No accessible keyboard devices found")
	}

	fmt.Println("Press keys anywhere. Ctrl+C to exit.")
	wg.Wait()
}
