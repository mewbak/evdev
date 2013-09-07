// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/evdev"
	"os"
	"os/signal"
	"strings"
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

	// Fetch the driver version for this device.
	major, minor, revision := dev.Version()

	// Fetch device identity.
	id := dev.Id()

	// Fetch the bitset, specifying the supported event types.
	events := dev.EventTypes()

	// Fetch device name.
	name := dev.Name()

	// Fetch physical device mapping.
	path := dev.Path()

	// Fetch device serial code.
	// For most devices this is not defined and yields an empty string.
	serial := dev.Serial()

	// Display all the collected information about our device.
	fmt.Printf(" Node    : %s\n", node)
	fmt.Printf(" Name    : %s\n", name)
	fmt.Printf(" Path    : %s\n", path)
	fmt.Printf(" Serial  : %s\n", serial)
	fmt.Printf(" Driver  : %d.%d.%d\n", major, minor, revision)
	fmt.Printf(" Vendor  : %04x\n", id.Vendor)
	fmt.Printf(" Product : %04x\n", id.Product)
	fmt.Printf(" Version : %04x\n", id.Version)
	fmt.Printf(" Bus     : %s\n", busName(id.BusType))
	fmt.Printf(" Events  : %s\n", listEvents(events))

	// Read events from the device, until we exit the program.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	for {
		select {
		case <-signals:
			return
		case evt := <-dev.Inbox:
			fmt.Printf("%+v\n", evt)
		}
	}
}

// busName returns the string equivalent of the given bus type.
func busName(bus uint16) string {
	switch bus {
	case evdev.BusPCI:
		return "PCI"
	case evdev.BusISAPNP:
		return "ISA Plug & Play"
	case evdev.BusUSB:
		return "USB"
	case evdev.BusHIL:
		return "HIL"
	case evdev.BusBluetooth:
		return "Bluetooth"
	case evdev.BusVirtual:
		return "Virtual"
	case evdev.BusISA:
		return "ISA"
	case evdev.BusI8042:
		return "I8042"
	case evdev.BusXTKBD:
		return "XT Keyboard"
	case evdev.BusRS232:
		return "RS232"
	case evdev.BusGamePort:
		return "Game Port"
	case evdev.BusParPort:
		return "Parallel Port"
	case evdev.BusAmiga:
		return "Amiga"
	case evdev.BusADB:
		return "ADB"
	case evdev.BusI2C:
		return "I2C"
	case evdev.BusHost:
		return "Host"
	case evdev.BusGSC:
		return "GSC"
	case evdev.BusAtari:
		return "Atari"
	case evdev.BusSPI:
		return "SPI"
	}

	return "Unknown"
}

// listEvents lists the event types supported by the device.
func listEvents(set evdev.Bitset) string {
	var list []string

	for n := 0; n < set.Len(); n++ {
		if !set.Test(n) {
			continue
		}

		switch n {
		case evdev.EvSync:
			list = append(list, "Sync Events")
		case evdev.EvKey:
			list = append(list, "Keys or Buttons")
		case evdev.EvRelative:
			list = append(list, "Relative Axes")
		case evdev.EvAbsolute:
			list = append(list, "Absolute Axes")
		case evdev.EvMisc:
			list = append(list, "Miscellaneous")
		case evdev.EvLed:
			list = append(list, "LEDs")
		case evdev.EvSound:
			list = append(list, "Sounds")
		case evdev.EvRepeat:
			list = append(list, "Repeat")
		case evdev.EvForceFeedback,
			evdev.EvForceFeedbackStatus:
			list = append(list, "Force Feedback")
		case evdev.EvPower:
			list = append(list, "Power Management")
		case evdev.EvSw:
			list = append(list, "Binary switches")
		default:
			list = append(list, fmt.Sprintf("Unknown (0x%02x)", n))
		}
	}

	return strings.Join(list, ", ")
}

func parseArgs() string {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <node>\n", os.Args[0])
		os.Exit(1)
	}

	return flag.Args()[0]
}
