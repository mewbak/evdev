// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import "unsafe"

// Relative events describe relative changes in a property.
// For example, a mouse may move to the left by a certain
// number of units, but its absolute position in space is unknown.
// If the absolute position is known, EvAbsolute codes should be used
// instead of EvRelative codes.
//
// RelWheel and RelHWheel are used for vertical and horizontal scroll
// wheels, respectively.
const (
	RelX      = 0x00
	RelY      = 0x01
	RelZ      = 0x02
	RelRX     = 0x03
	RelRY     = 0x04
	RelRZ     = 0x05
	RelHWheel = 0x06
	RelDial   = 0x07
	RelWheel  = 0x08
	RelMisc   = 0x09
	RelMax    = 0x0f
	RelCount  = RelMax + 1
)

// RelativeAxes returns a bitfield indicating which relative axes are
// supported by the device.
//
// This is only applicable to devices with EvRelative event support.
func (d *Device) RelativeAxes() Bitset {
	bs := NewBitset(RelMax)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGBIT(EvRelative, len(buf)), unsafe.Pointer(&buf[0]))
	return bs
}
