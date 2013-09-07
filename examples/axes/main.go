// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/evdev"
	"os"
)

func main() {
	node := parseArgs()

	// Create and open our device.
	dev, err := evdev.Open(node)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// Make sure it is closed once we are done.
	defer dev.Close()

	// Ensure this device supports the needed event types.
	if !hasType(dev, evdev.EvAbsolute) {
		fmt.Fprintf(os.Stderr, "Device %q does not support absolute axis events.\n", node)
		return
	}

	// Fetch the supported axes.
	axes := dev.Axes()
	for n := 0; n < axes.Len(); n++ {
		if !axes.Test(n) {
			continue
		}

		fmt.Printf("  Absolute axis 0x%02x ", n)

		switch n {
		case evdev.AbsX:
			fmt.Printf("X Axis: ")
		case evdev.AbsY:
			fmt.Printf("Y Axis: ")
		case evdev.AbsZ:
			fmt.Printf("Z Axis: ")
		default: // More axes types...
			fmt.Printf("Other axis\n")
		}

		// Get axis information.
		abs := dev.AbsInfo(n)
		fmt.Printf("%+v\n", abs)
	}
}

// hasKeys determines if the given device supports the specified
// event type. E.g.: EvKey, EvAbsolute, etc.
func hasType(dev *evdev.Device, etype int) bool {
	events := dev.EventTypes()

	for n := 0; n < events.Len(); n++ {
		if n == etype && events.Test(n) {
			return true
		}
	}

	return false
}

func parseArgs() string {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <node>\n", os.Args[0])
		os.Exit(1)
	}

	return flag.Args()[0]
}
