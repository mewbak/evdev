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

	// Ensure proper cleanup when things go booboo.
	defer func() {
		if err != nil {
			dev.Close()
		}
	}()

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

// Version returns version information for the device driver.
// These being major, minor and revision numbers.
func (d *Device) Version() (int, int, int, error) {
	var version uint32
	err := ioctl(d.fd.Fd(), _EVIOCGVERSION, unsafe.Pointer(&version))
	if err != nil {
		return 0, 0, 0, err
	}

	return int(version>>16) & 0xffff, int(version>>8) & 0xff,
		int(version) & 0xff, nil
}

// Id returns identification information for the given device.
func (d *Device) Id() (Id, error) {
	var id Id
	err := ioctl(d.fd.Fd(), _EVIOCGID, unsafe.Pointer(&id))
	return id, err
}
