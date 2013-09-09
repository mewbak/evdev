// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"errors"
	"fmt"
	"os"
)

// List of device types.
//
// These are used to look for specific input device types
// using evdev.Find().
//
// The returned devices may not necessarily be an actual
// keyboard or mouse, etc. Just a device which can behave like one.
// For instance: Mouse may return a trackpad, a multi-touch screen
// and an actual mouse if all of these happen to be connected.
// It is up to the host to figure out which one to use.
const (
	Keyboard = iota
	Mouse
	Joystick
)

// Find returns a list of all attached devices, which
// qualify as the given device type.
func Find(devtype int) (list []*Device, err error) {
	// Ensure we clean up properly if something goes wrong.
	defer func() {
		if err != nil {
			for _, dev := range list {
				dev.Close()
			}
			list = nil
		}
	}()

	var count int
	var dev *Device
	var testFunc func(*Device) bool

	switch devtype {
	case Keyboard:
		testFunc = IsKeyboard
	case Mouse:
		testFunc = IsMouse
	case Joystick:
		testFunc = IsJoystick
	default:
		err = errors.New("Invalid device type")
		return
	}

	for {
		node := fmt.Sprintf("/dev/input/event%d", count)
		dev, err = Open(node)

		if err != nil {
			if os.IsNotExist(err) {
				err = nil
			}
			return
		}

		count++

		if testFunc(dev) {
			list = append(list, dev)
		}
	}

	return
}

// IsKeyboard returns true if the given device qualifies as a keyboard.
func IsKeyboard(dev *Device) bool {
	return dev.Test(dev.EventTypes(), EvKeys, EvLed)
}

// IsMouse returns true if the given device qualifies as a mouse.
func IsMouse(dev *Device) bool {
	return dev.Test(dev.EventTypes(), EvKeys, EvRelative)
}

// IsJoystick returns true if the given device qualifies as a joystick.
func IsJoystick(dev *Device) bool {
	return dev.Test(dev.EventTypes(), EvKeys, EvAbsolute)
}
