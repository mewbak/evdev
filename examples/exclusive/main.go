// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/evdev"
	"os"
	"time"
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

	// Obtain an exclusive lock on the device.
	if !dev.Grab() {
		fmt.Fprintf(os.Stderr, "Failed to abtain device lock.\n")
		return
	}

	fmt.Printf("%s is now locked\n", dev.Name())

	// Do work with device...
	<-time.After(5e9)

	// Release lock
	if !dev.Release() {
		fmt.Fprintf(os.Stderr, "Failed to release device lock.\n")
		return
	}

	fmt.Printf("%s is now unlocked\n", dev.Name())
}

func parseArgs() string {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <node>\n", os.Args[0])
		os.Exit(1)
	}

	return flag.Args()[0]
}
