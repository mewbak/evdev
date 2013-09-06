// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"fmt"
	"github.com/jteeuwen/evdev"
	"os"
)

const Device = "/dev/input/event0"

func main() {
	// Create and open our device.
	dev, err := evdev.Open(Device)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// Make sure it is closed once we are done.
	defer dev.Close()

	// Fetch the driver version for this device.
	major, minor, revision, err := dev.Version()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	fmt.Printf("Driver version: %d.%d.%d\n", major, minor, revision)

	// Fetch device id attributes.
	id, err := dev.Id()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	fmt.Printf("Id: %+v\n", id)
}
