// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"syscall"
	"unsafe"
)

// <linux/input.h>

const EvVersion = 0x010001

// Event represents a generic input event.
type Event struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

// Event types
const (
	EvSync                = 0x00 // Synchronisation events.
	EvKeys                = 0x01 // Absolute binary results, such as keys and buttons.
	EvRelative            = 0x02 // Relative results, such as the axes on a mouse.
	EvAbsolute            = 0x03 // Absolute integer results, such as the axes on a joystick or for a tablet.
	EvMisc                = 0x04 // Miscellaneous uses that didn't fit anywhere else.
	EvSwitch              = 0x05 // Used to describe stateful binary switches.
	EvLed                 = 0x11 // LEDs and similar indications.
	EvSound               = 0x12 // Sound output, such as buzzers.
	EvRepeat              = 0x14 // Enables autorepeat of keys in the input core.
	EvForceFeedback       = 0x15 // Sends force-feedback effects to a device.
	EvPower               = 0x16 // Power management events.
	EvForceFeedbackStatus = 0x17 // Device reporting of force-feedback effects back to the host.
	EvMax                 = 0x1f
	EvCount               = EvMax + 1
)

// EventTypes determines the device's capabilities.
// It yields a bitset which can be tested against
// EvXXX constants to determine which types are supported.
func (d *Device) EventTypes() Bitset {
	bs := NewBitset(EvMax)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGBIT(0, EvMax), unsafe.Pointer(&buf[0]))
	return bs
}

// IDs.
const (
	IdBus = iota
	IdVendor
	IdProduct
	IdVersion
)
