// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"syscall"
	"unsafe"
)

// <linux/input.h>

const EvVersion = 0x010001

// Event types
const (
	EvSync                = 0x00 // Synchronisation events.
	EvKeys                = 0x01 // Absolute binary results, such as keys and buttons.
	EvRelative            = 0x02 // Relative results, such as the axes on a mouse.
	EvAbsolute            = 0x03 // Absolute integer results, such as the axes on a joystick or for a tablet.
	EvMisc                = 0x04 // Miscellaneous uses that didn't fit anywhere else.
	EvSwitch              = 0x05 // Used to describe binary state input switches
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

var (
	_EVIOCGVERSION    uintptr
	_EVIOCGID         uintptr
	_EVIOCGREP        uintptr
	_EVIOCSREP        uintptr
	_EVIOCGKEYCODE    uintptr
	_EVIOCGKEYCODE_V2 uintptr
	_EVIOCSKEYCODE    uintptr
	_EVIOCSKEYCODE_V2 uintptr
	_EVIOCSFF         uintptr
	_EVIOCRMFF        uintptr
	_EVIOCGEFFECTS    uintptr
	_EVIOCGRAB        uintptr
	_EVIOCSCLOCKID    uintptr
)

func init() {
	var i int32
	var id Id
	var ke KeymapEntry
	var ffe Effect

	sizeof_int := int(unsafe.Sizeof(i))
	sizeof_int2 := sizeof_int << 1
	sizeof_id := int(unsafe.Sizeof(id))
	sizeof_keymap_entry := int(unsafe.Sizeof(ke))
	sizeof_effect := int(unsafe.Sizeof(ffe))

	_EVIOCGVERSION = uintptr(_IOR('E', 0x01, sizeof_int))
	_EVIOCGID = uintptr(_IOR('E', 0x02, sizeof_id))
	_EVIOCGREP = uintptr(_IOR('E', 0x03, sizeof_int2))
	_EVIOCSREP = uintptr(_IOW('E', 0x03, sizeof_int2))

	_EVIOCGKEYCODE = uintptr(_IOR('E', 0x04, sizeof_int2))
	_EVIOCGKEYCODE_V2 = uintptr(_IOR('E', 0x04, sizeof_keymap_entry))
	_EVIOCSKEYCODE = uintptr(_IOW('E', 0x04, sizeof_int2))
	_EVIOCSKEYCODE_V2 = uintptr(_IOW('E', 0x04, sizeof_keymap_entry))

	_EVIOCSFF = uintptr(_IOC(_IOC_WRITE, 'E', 0x80, sizeof_effect))
	_EVIOCRMFF = uintptr(_IOW('E', 0x81, sizeof_int))
	_EVIOCGEFFECTS = uintptr(_IOR('E', 0x84, sizeof_int))
	_EVIOCGRAB = uintptr(_IOW('E', 0x90, sizeof_int))
	_EVIOCSCLOCKID = uintptr(_IOW('E', 0xa0, sizeof_int))

}

func _EVIOCGNAME(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x06, len)
}

func _EVIOCGPHYS(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x07, len)
}

func _EVIOCGUNIQ(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x08, len)
}

func _EVIOCGPROP(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x09, len)
}

func _EVIOCGMTSLOTS(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x0a, len)
}

func _EVIOCGKEY(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x18, len)
}

func _EVIOCGLED(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x19, len)
}

func _EVIOCGSND(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x1a, len)
}

func _EVIOCGSW(len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x1b, len)
}

func _EVIOCGBIT(ev, len int) uintptr {
	return _IOC(_IOC_READ, 'E', 0x20+ev, len)
}

func _EVIOCGABS(abs int) uintptr {
	var v AbsInfo
	return _IOR('E', 0x40+abs, int(unsafe.Sizeof(v)))
}

func _EVIOCSABS(abs int) uintptr {
	var v AbsInfo
	return _IOW('E', 0xc0+abs, int(unsafe.Sizeof(v)))
}

// Event represents a generic input event.
type Event struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

// Device properties and quirks
const (
	InputPropPointer   = 0x00 // needs a pointer
	InputPropDirect    = 0x01 // direct input devices
	InputPropButtonPad = 0x02 // has button(s) under pad
	InputPropSemiMT    = 0x03 // touch rectangle only
	InputPropMax       = 0x1f
	InputPropCount     = InputPropMax + 1
)

// Synchronization events.
const (
	SynReport = iota
	SynConfig
	SynMTReport
	SynDropped
)

// Misc events
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

// MTTool types
const (
	MtToolFinger = 0
	MtToolPen    = 1
	MtToolMax    = 1
)

// IDs.
const (
	IdBus = iota
	IdVendor
	IdProduct
	IdVersion
)
