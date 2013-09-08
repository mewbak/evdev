// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import "unsafe"

// LEDs
const (
	LedNumLock    = 0x00
	LedCapsLock   = 0x01
	LedScrollLock = 0x02
	LedCompose    = 0x03
	LedKana       = 0x04
	LedSleep      = 0x05
	LedSuspend    = 0x06
	LedMute       = 0x07
	LedMisc       = 0x08
	LedMail       = 0x09
	LedCharging   = 0x0a
	LedMax        = 0x0f
	LedCount      = LedMax + 1
)

// LEDState returns the current, global LED state.
//
// This is only applicable to devices with EvLed event support.
func (d *Device) LEDState() Bitset {
	bs := NewBitset(LedMax)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGLED(len(buf)), unsafe.Pointer(&buf[0]))
	return bs
}
