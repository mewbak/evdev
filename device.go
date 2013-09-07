// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"fmt"
	"os"
	"sync"
	"unsafe"
)

const eventBufferSize = 64

// Device represents a single device node.
type Device struct {
	wg     sync.WaitGroup
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

	dev.wg.Add(2)
	go dev.pollIn()
	go dev.pollOut()
	return
}

// Close closes the underlying device node.
func (d *Device) Close() (err error) {
	if d.fd != nil {
		err = d.fd.Close()
		d.fd = nil
	}

	close(d.Inbox)
	close(d.Outbox)
	d.wg.Wait()
	return
}

// Keystate returns the current,global key- and button- states.
// This applies only to devices which have the EvKey capability defined.
// This can be determined through `Device.EventTypes()`.
func (d *Device) KeyState() Bitset {
	bs := NewBitset(KeyMax)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGKEY(len(buf)), unsafe.Pointer(&buf[0]))
	return bs
}

// LEDState returns the current, global LED state.
// This applies only to devices which have the EvLed capability defined.
// This can be determined through `Device.EventTypes()`.
func (d *Device) LEDState() Bitset {
	bs := NewBitset(LedMax)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGLED(len(buf)), unsafe.Pointer(&buf[0]))
	return bs
}

// KeyMap fills the key mapping for the given key.
// E.g.: Pressing M, will input N into the input system.
// This allows us to rewire physical keys.
//
// Refer to `Device.SetKeyMap()` for information on what
// this means.
//
// Be aware that the KeyMap functions may not work on every keyboard.
func (d *Device) KeyMap(key int) int {
	var m [2]int32
	m[0] = int32(key)
	ioctl(d.fd.Fd(), _EVIOCGKEYCODE, unsafe.Pointer(&m[0]))
	return int(m[1])
}

// SetKeyMap sets the given key to the specified mapping.
// E.g.: Pressing M, will input N into the input system.
// This allows us to rewire physical keys.
//
// Some input drivers support variable mappings between the keys
// held down (which are interpreted by the keyboard scan and reported
// as scancodes) and the events sent to the input layer.
//
// You can change which key is associated with each scancode
// using this call. The value of the scancode is the first element
// in the integer array (list[n][0]), and the resulting input
// event key number (keycode) is the second element in the array.
// (list[n][1]).
//
// Be aware that the KeyMap functions may not work on every keyboard.
func (d *Device) SetKeyMap(key, value int) {
	var m [2]int32
	m[0] = int32(key)
	m[1] = int32(value)
	ioctl(d.fd.Fd(), _EVIOCSKEYCODE, unsafe.Pointer(&m[0]))
}

// RepeatState returns the current, global repeat state.
// This applies only to devices which have the EvRepeat capability defined.
// This can be determined through `Device.EventTypes()`.
//
// Refer to Device.SetRepeatState for an explanation on what
// the returned values mean.
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
func (d *Device) SetRepeatState(initial, subsequent uint) bool {
	var rep [2]int32
	rep[0] = int32(initial)
	rep[1] = int32(subsequent)
	return ioctl(d.fd.Fd(), _EVIOCSREP, unsafe.Pointer(&rep[0])) == nil
}

// Axes returns a bitfield indicating which axes are
// supported by the device.
func (d *Device) Axes() Bitset {
	bs := NewBitset(AbsMax)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGBIT(EvAbsolute, len(buf)), unsafe.Pointer(&buf[0]))
	return bs
}

// AbsInfo provides state information for one absolute axis.
// If you want the global state for a device, you have to call
// the function for each axis present on the device.
// See Device.Axes() for details on how find them.
func (d *Device) AbsInfo(axis int) AbsInfo {
	var abs AbsInfo
	ioctl(d.fd.Fd(), _EVIOCGABS(axis), unsafe.Pointer(&abs))
	return abs
}

// ForceFeedbackCaps returns a bitset which specified the kind of Force Feedback
// effects supported by this device. The bits can be compared against
// the FFXXX constants. Additionally, it returns the number of effects
// this device can handle simultaneously.
//
// This is only applicable to devices with the EvForceFeedback event type set.
func (d *Device) ForceFeedbackCaps() (int, Bitset) {
	bs := NewBitset(24)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGBIT(EvForceFeedback, len(buf)), unsafe.Pointer(&buf[0]))

	var count int32
	ioctl(d.fd.Fd(), _EVIOCGEFFECTS, unsafe.Pointer(&count))
	return int(count), bs
}

// Supports determines if the device supports the specified
// event types. E.g.: EvKey, EvAbsolute, etc.
//
// It returns true only if the device supports all the supplied
// event types.
func (d *Device) Supports(etype ...int) bool {
	events := d.EventTypes()

	for i := range etype {
		for n := 0; n < events.Len(); n++ {
			if n == etype[i] && events.Test(n) {
				etype[i] = -1
			}
		}
	}

	for i := range etype {
		if etype[i] > 0 {
			return false
		}
	}

	return true
}

// EventTypes determines the device's capabilities.
// It yields a bitset which can be tested against
// EvXXX constants to determine which types are supported.
func (d *Device) EventTypes() Bitset {
	bs := NewBitset(EvMax)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGBIT(0, EvMax), unsafe.Pointer(&buf[0]))
	return bs
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
	ioctl(d.fd.Fd(), _EVIOCGPHYS(256), unsafe.Pointer(&str[0]))
	return string(str[:])
}

// Serial returns the unique serial code for the device.
// Most devices do not have this and will return an empty string.
func (d *Device) Serial() string {
	var str [256]byte
	ioctl(d.fd.Fd(), _EVIOCGUNIQ(256), unsafe.Pointer(&str[0]))
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
	defer d.wg.Done()

	var e Event

	size := int(unsafe.Sizeof(e))
	buf := make([]byte, size*eventBufferSize)

	for {
		n, err := d.fd.Read(buf)
		if err != nil {
			return
		}

		evt := (*(*[1<<31 - 1]Event)(unsafe.Pointer(&buf[0])))[:n/size]

		for n = range evt {
			d.Inbox <- evt[n]
		}
	}
}

// pollOut polls the outbox for pending messages.
// These are then sent to the device.
func (d *Device) pollOut() {
	defer d.wg.Done()

	var e Event
	size := int(unsafe.Sizeof(e))

	for msg := range d.Outbox {
		buf := (*(*[1<<31 - 1]byte)(unsafe.Pointer(&msg)))[:size]

		n, err := d.fd.Write(buf)
		if err != nil {
			return
		}

		if n < size {
			fmt.Fprintf(os.Stderr, "poll outbox: short write\n")
		}
	}
}
