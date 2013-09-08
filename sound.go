// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

// Sound events are used for sending sound
// commands to simple sound output devices.
const (
	SndClick = 0x00
	SndBell  = 0x01
	SndTone  = 0x02
	SndMax   = 0x07
	SndCount = SndMax + 1
)
