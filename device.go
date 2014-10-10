// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"fmt"
	"os"
	"unsafe"
)

const eventBufferSize = 64

// Device represents a single device node.
type Device struct {
	fd     *os.File
	Inbox  chan Event // Channel exposing incoming events.
	Outbox chan Event // Channel for outgoing events.
}

// Open opens a new device for the given node name.
// This can be anything listed in /dev/input/event[x].
func Open(node string) (dev *Device, err error) {
	dev = new(Device)
	dev.fd, err = os.OpenFile(node, os.O_RDWR, 0)
	dev.Inbox = make(chan Event, eventBufferSize)
	dev.Outbox = make(chan Event, 1)

	go dev.pollIn()
	go dev.pollOut()
	return
}

// Close closes the underlying device node.
func (d *Device) Close() (err error) {
	d.Release()

	if d.fd != nil {
		err = d.fd.Close()
		d.fd = nil
	}

	return
}

// Grab attempts to gain exclusive access to this device.
// This means that we are the only ones receiving events from
// the device; other processes will not.
//
// This ability should be handled with care, especially
// when trying to lock keyboard access. If this is
// executed while we are running in something like X,
// this call will prevent X from receiving any and all
// keyboard events. All of them will only be sent to our
// own process. If we do not properly handle these key
// events, we may lock ourselves out of the system
// and a hard reset is required to restore it.
func (d *Device) Grab() bool {
	return ioctl(d.fd.Fd(), _EVIOCGRAB, 1) == nil
}

// Release releases a lock, previously obtained through `Device.Grab`.
func (d *Device) Release() bool {
	return ioctl(d.fd.Fd(), _EVIOCGRAB, 0) == nil
}

// Test takes a bitset and a list of constants
// and tests one against the other to see if the device
// supports a given set of properties.
//
// It returns true only if the bitset defines all the supplied types.
// E.g.: To test for certain event types:
//
//     if dev.Test(dev.EventTypes(), EvKey, EvRepeat) {
//
// To test for certain absolute axes:
//
//     if dev.Test(dev.AbsoluteAxes(), AbsX, AbsY, AbsZ) {
//
// To test for certain relative axes:
//
//     if dev.Test(dev.RelativeAxes(), RelX, RelY, RelZ) {
func (d *Device) Test(set Bitset, values ...int) bool {
	var count int

	for i := range values {
		for n := 0; n < set.Len(); n++ {
			if n == values[i] && set.Test(n) {
				count++
				break
			}
		}
	}

	return count == len(values)
}

// Name returns the name of the device.
func (d *Device) Name() string {
	var str [256]byte
	ioctl(d.fd.Fd(), _EVIOCGNAME(256), unsafe.Pointer(&str[0]))
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
	ioctl(d.fd.Fd(), _EVIOCGPHYS(len(str)), unsafe.Pointer(&str[0]))
	return string(str[:])
}

// Serial returns the unique serial code for the device.
// Most devices do not have this and will return an empty string.
func (d *Device) Serial() string {
	var str [256]byte
	ioctl(d.fd.Fd(), _EVIOCGUNIQ(len(str)), unsafe.Pointer(&str[0]))
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

// pollIn polls the device for incoming events.
// We can receive many events with a single read.
// This is why the outgoing event channel has a large buffer.
func (d *Device) pollIn() {
	defer close(d.Inbox)

	var e Event

	size := int(unsafe.Sizeof(e))
	buf := make([]byte, size*eventBufferSize)

	for {
		n, err := d.fd.Read(buf)
		if err != nil {
			return
		}

		evt := (*(*[1<<27 - 1]Event)(unsafe.Pointer(&buf[0])))[:n/size]

		for n = range evt {
			d.Inbox <- evt[n]
		}
	}
}

// pollOut polls the outbox for pending messages.
// These are then sent to the device.
func (d *Device) pollOut() {
	defer close(d.Outbox)

	var e Event
	size := int(unsafe.Sizeof(e))

	for msg := range d.Outbox {
		buf := (*(*[1<<27 - 1]byte)(unsafe.Pointer(&msg)))[:size]

		n, err := d.fd.Write(buf)
		if err != nil {
			return
		}

		if n < size {
			fmt.Fprintf(os.Stderr, "poll outbox: short write\n")
		}
	}
}
