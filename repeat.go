// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import "unsafe"

// Autorepeat values
const (
	RepDelay  = 0x00
	RepPeriod = 0x01
	RepMax    = 0x01
	RepCount  = RepMax + 1
)

// RepeatState returns the current, global repeat state.
// This applies only to devices which have the EvRepeat capability defined.
// This can be determined through `Device.EventTypes()`.
//
// Refer to Device.SetRepeatState for an explanation on what
// the returned values mean.
//
// This is only applicable to devices with EvRepeat event support.
func (d *Device) RepeatState() (uint, uint) {
	var rep [2]int32
	ioctl(d.fd.Fd(), _EVIOCGREP, unsafe.Pointer(&rep[0]))
	return uint(rep[0]), uint(rep[1])
}

// SetRepeatState sets the global repeat state for the given
// device.
//
// The values indicate (in milliseconds) the delay before
// the device starts repeating and the delay between
// subsequent repeats. This might apply to a keyboard where
// the user presses and holds a key.
//
// E.g.: We see an initial character immediately, then
// another @initial milliseconds later and after that,
// once every @subsequent milliseconds, until the key
// is released.
//
// This returns false if the operation failed.
//
// This is only applicable to devices with EvRepeat event support.
func (d *Device) SetRepeatState(initial, subsequent uint) bool {
	var rep [2]int32
	rep[0] = int32(initial)
	rep[1] = int32(subsequent)
	return ioctl(d.fd.Fd(), _EVIOCSREP, unsafe.Pointer(&rep[0])) == nil
}
