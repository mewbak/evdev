// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/evdev"
	"os"
	"os/signal"
)

var (
	keyboards []*evdev.Device
	mice      []*evdev.Device
)

func main() {
	var err error

	defer cleanup()

	keyboards, err = evdev.Find(evdev.Keyboard)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Find keyboards: %v\n", err)
		return
	}

	mice, err = evdev.Find(evdev.Mouse)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Find mice: %v\n", err)
		return
	}

	fmt.Printf("Keyboards: %d\n", len(keyboards))
	for _, dev := range keyboards {
		fmt.Printf(" - %s at %s\n", dev.Name(), dev.Path())
	}

	fmt.Printf("Mice: %d\n", len(keyboards))
	for _, dev := range mice {
		fmt.Printf(" - %s at %s\n", dev.Name(), dev.Path())
	}

	if len(keyboards) == 0 || len(mice) == 0 {
		return
	}

	// Poll for events or exit signals.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	for {
		select {
		case <-signals:
			return

		case evt := <-keyboards[0].Inbox:
			fmt.Printf("Keyboard: %v\n", evt)

		case evt := <-mice[0].Inbox:
			fmt.Printf("Mouse: %v\n", evt)
		}
	}
}

func cleanup() {
	for _, dev := range keyboards {
		dev.Close()
	}

	for _, dev := range mice {
		dev.Close()
	}
}
