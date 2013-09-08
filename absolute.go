// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import "unsafe"

// Absolute axes
const (
	AbsX             = 0x00
	AbsY             = 0x01
	AbsZ             = 0x02
	AbsRX            = 0x03
	AbsRY            = 0x04
	AbsRZ            = 0x05
	AbsThrottle      = 0x06
	AbsRudder        = 0x07
	AbsWheel         = 0x08
	AbsGas           = 0x09
	AbsBrake         = 0x0a
	AbsHat0X         = 0x10
	AbsHat0Y         = 0x11
	AbsHat1X         = 0x12
	AbsHat1Y         = 0x13
	AbsHat2X         = 0x14
	AbsHat2Y         = 0x15
	AbsHat3X         = 0x16
	AbsHat3Y         = 0x17
	AbsPressure      = 0x18
	AbsDistance      = 0x19
	AbsTiltX         = 0x1a
	AbsTiltY         = 0x1b
	AbsToolWidth     = 0x1c
	AbsVolume        = 0x20
	AbsMisc          = 0x28
	AbsMTSlot        = 0x2f // MT slot being modified
	AbsMTTouchMajor  = 0x30 // Major axis of touching ellipse
	AbsMTTouchMinor  = 0x31 // Minor axis (omit if circular)
	AbsMTWidthMajor  = 0x32 // Major axis of approaching ellipse
	AbsMTWidthMinor  = 0x33 // Minor axis (omit if circular)
	AbsMTOrientation = 0x34 // Ellipse orientation
	AbsMTPositionX   = 0x35 // Center X touch position
	AbsMTPositionY   = 0x36 // Center Y touch position
	AbsMTToolTYPE    = 0x37 // Type of touching device
	AbsMTBlobId      = 0x38 // Group a set of packets as a blob
	AbsMTTrackingId  = 0x39 // Unique ID of initiated contact
	AbsMTPressure    = 0x3a // Pressure on contact area
	AbsMTDistance    = 0x3b // Contact hover distance
	AbsMTToolX       = 0x3c // Center X tool position
	AbsMTToolY       = 0x3d // Center Y tool position
	AbsMax           = 0x3f
	AbsCount         = AbsMax + 1
)

// AbsInfo provides information for a specific absolute axis.
// This applies to devices which support EvAbsolute events.
type AbsInfo struct {
	Value      int32 // Current value of the axis,
	Minimum    int32 // Lower limit of axis.
	Maximum    int32 // Upper limit of axis.
	Fuzz       int32 // ???
	Flat       int32 // Size of the 'flat' section.
	Resolution int32 // Size of the error that may be present.
}

// AbsoluteAxes returns a bitfield indicating which absolute axes are
// supported by the device.
//
// This is only applicable to devices with EvAbsolute event support.
func (d *Device) AbsoluteAxes() Bitset {
	bs := NewBitset(AbsMax)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGBIT(EvAbsolute, len(buf)), unsafe.Pointer(&buf[0]))
	return bs
}

// AbsoluteInfo provides state information for one absolute axis.
// If you want the global state for a device, you have to call
// the function for each axis present on the device.
// See Device.AbsoluteAxes() for details on how find them.
//
// This is only applicable to devices with EvAbsolute event support.
func (d *Device) AbsoluteInfo(axis int) AbsInfo {
	var abs AbsInfo
	ioctl(d.fd.Fd(), _EVIOCGABS(axis), unsafe.Pointer(&abs))
	return abs
}
