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

	// Ensure this device supports key/button events.
	if !dev.SupportsEvent(evdev.EvKey) {
		fmt.Fprintf(os.Stderr, "Device %q does not support key/button events.\n", node)
		return
	}

	// Fetch the current key/button state and display it.
	ks := dev.KeyState()
	listState(ks)
}

// listState prints the global key/button state, as defined
// in the given bitset.
func listState(set evdev.Bitset) {
	for n := 0; n < set.Len(); n++ {
		// The key is considered pressed if the bitset
		// has its corresponding bit set.
		if !set.Test(n) {
			continue
		}

		fmt.Printf("  Key 0x%02x ", n)

		switch n {
		case evdev.KeyReserved:
			fmt.Printf("Reserved")
		case evdev.KeyEscape:
			fmt.Printf("Escape")
		case evdev.BtnStylus2:
			fmt.Printf("2nd stylus button")

			// more keys/buttons..
		}

		fmt.Println()
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
