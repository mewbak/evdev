// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"os"
	"unsafe"
)

// Device represents a single input device node.
type Device struct {
	fd *os.File
}

// Open opens a new device for the given node name.
// This can be anything listed in /dev/input/.
func Open(node string) (dev *Device, err error) {
	dev = new(Device)
	dev.fd, err = os.Open(node)
	return
}

// Close closes the underlying device node.
func (d *Device) Close() (err error) {
	if d.fd != nil {
		err = d.fd.Close()
		d.fd = nil
	}
	return
}

// EventTypes determines the device's capabilities.
// It yields a bit mask which can be tested for against
// EvXXX constants to determine which types are supported.
func (d *Device) EventTypes() uint64 {
	var bitmask uint64
	ioctl(d.fd.Fd(), uintptr(_EVIOCGBIT(0, 8)), unsafe.Pointer(&bitmask))
	return bitmask
}

// Name returns the name of the device.
func (d *Device) Name() string {
	var str [256]byte
	ioctl(d.fd.Fd(), uintptr(_EVIOCGNAME(256)), unsafe.Pointer(&str[0]))
	return string(str[:])
}

// Path returns the physical path of the device.
// For example:
//
//    usb-00:01.2-2.1/input0
//
// To understand what this string is showing, you need
// to break it down into parts. `usb` means this is
// a physical topology from the USB system.
//
// `00:01.2` is the PCI bus information for the USB host
// controller (in this case, bus 0, slot 1, function 2).
//
// `2.1` shows the path from the root hub to the device.
// In this case, the upstream hub is plugged in to the
// second port on the root hub, and that device is plugged
// in to the first port on the upstream hub.
//
// `input0` means this is the first event interface on the device.
// Most devices have only one, but multimedia keyboards
// may present the normal keyboard on one interface and
// the multimedia function keys on a second interface.
func (d *Device) Path() string {
	var str [256]byte
	ioctl(d.fd.Fd(), uintptr(_EVIOCGPHYS(256)), unsafe.Pointer(&str[0]))
	return string(str[:])
}

// Serial returns the unique serial code for the device.
// Most devices do not have this and will return an empty string.
func (d *Device) Serial() string {
	var str [256]byte
	ioctl(d.fd.Fd(), uintptr(_EVIOCGUNIQ(256)), unsafe.Pointer(&str[0]))
	return string(str[:])
}

// Version returns version information for the device driver.
// These being major, minor and revision numbers.
func (d *Device) Version() (int, int, int) {
	var version uint32
	err := ioctl(d.fd.Fd(), _EVIOCGVERSION, unsafe.Pointer(&version))
	if err != nil {
		return 0, 0, 0
	}

	return int(version>>16) & 0xffff,
		int(version>>8) & 0xff,
		int(version) & 0xff
}

// Id returns the device identity.
func (d *Device) Id() Id {
	var id Id
	ioctl(d.fd.Fd(), _EVIOCGID, unsafe.Pointer(&id))
	return id
}
