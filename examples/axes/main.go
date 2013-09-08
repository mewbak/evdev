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

	events := dev.EventTypes()
	abs := dev.Test(events, evdev.EvAbsolute)
	rel := dev.Test(events, evdev.EvRelative)

	if !abs && !rel {
		fmt.Fprintf(os.Stderr, "Device %q does not support relative or absolute axes.\n", node)
		return
	}

	if abs {
		absAxes(dev)
	}

	if rel {
		absAxes(dev)
	}
}

// Testing for support of specific axes can be
// done with the `dev.Supports()` method.
func relAxes(dev *evdev.Device) {
	axes := dev.RelativeAxes()
	for n := 0; n < axes.Len(); n++ {
		if !axes.Test(n) {
			continue
		}

		fmt.Printf("  Relative axis 0x%02x ", n)

		switch n {
		case evdev.RelX:
			fmt.Printf("X Axis: ")
		case evdev.RelY:
			fmt.Printf("Y Axis: ")
		case evdev.RelZ:
			fmt.Printf("Z Axis: ")
		default: // More axes types...
			fmt.Printf("Other axis\n")
		}

		fmt.Println()
	}
}

// Testing for support of specific axes can be
// done with the `dev.Supports()` method.
func absAxes(dev *evdev.Device) {
	axes := dev.AbsoluteAxes()
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
		abs := dev.AbsoluteInfo(n)
		fmt.Printf("%+v\n", abs)
	}
}

func parseArgs() string {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <node>\n", os.Args[0])
		os.Exit(1)
	}

	return flag.Args()[0]
}
