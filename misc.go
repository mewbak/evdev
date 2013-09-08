// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

// Miscellaneous events are used for input and output events
// that do not fall under other categories.
//
// MiscTimestamp has a special meaning.
// It is used to report the number of microseconds since the last reset. This event
// should be coded as an uint32 value, which is allowed to wrap around with
// no special consequence. It is assumed that the time difference between two
// consecutive events is reliable on a reasonable time scale (hours).
// A reset to zero can happen, in which case the time since the last event is
// unknown. If the device does not provide this information, the driver must
// not provide it to user space.
const (
	MiscSerial    = 0x00
	MiscPulseLed  = 0x01
	MiscGesture   = 0x02
	MiscRaw       = 0x03
	MiscScan      = 0x04
	MiscTimestamp = 0x05
	MiscMax       = 0x07
	MiscCount     = MiscMax + 1
)
