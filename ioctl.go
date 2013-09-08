// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"fmt"
	"syscall"
	"unsafe"
)

func ioctl(fd, name uintptr, data interface{}) error {
	var v uintptr

	switch dd := data.(type) {
	case unsafe.Pointer:
		v = uintptr(dd)

	case int:
		v = uintptr(dd)

	case uintptr:
		v = dd

	default:
		return fmt.Errorf("ioctl: Invalid argument: %T", data)
	}

	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, v)
	if errno == 0 {
		return nil
	}

	return errno
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

const (
	_IOC_NONE      = 0x0
	_IOC_WRITE     = 0x1
	_IOC_READ      = 0x2
	_IOC_NRBITS    = 8
	_IOC_TYPEBITS  = 8
	_IOC_SIZEBITS  = 14
	_IOC_DIRBITS   = 2
	_IOC_NRSHIFT   = 0
	_IOC_NRMASK    = (1 << _IOC_NRBITS) - 1
	_IOC_TYPEMASK  = (1 << _IOC_TYPEBITS) - 1
	_IOC_SIZEMASK  = (1 << _IOC_SIZEBITS) - 1
	_IOC_DIRMASK   = (1 << _IOC_DIRBITS) - 1
	_IOC_TYPESHIFT = _IOC_NRSHIFT + _IOC_NRBITS
	_IOC_SIZESHIFT = _IOC_TYPESHIFT + _IOC_TYPEBITS
	_IOC_DIRSHIFT  = _IOC_SIZESHIFT + _IOC_SIZEBITS
	_IOC_IN        = _IOC_WRITE << _IOC_DIRSHIFT
	_IOC_OUT       = _IOC_READ << _IOC_DIRSHIFT
	_IOC_INOUT     = (_IOC_WRITE | _IOC_READ) << _IOC_DIRSHIFT
	_IOCSIZE_MASK  = _IOC_SIZEMASK << _IOC_SIZESHIFT
)

func _IOC(dir, t, nr, size int) uintptr {
	return uintptr((dir << _IOC_DIRSHIFT) | (t << _IOC_TYPESHIFT) |
		(nr << _IOC_NRSHIFT) | (size << _IOC_SIZESHIFT))
}

func _IO(t, nr int) uintptr {
	return _IOC(_IOC_NONE, t, nr, 0)
}

func _IOR(t, nr, size int) uintptr {
	return _IOC(_IOC_READ, t, nr, size)
}

func _IOW(t, nr, size int) uintptr {
	return _IOC(_IOC_WRITE, t, nr, size)
}

func _IOWR(t, nr, size int) uintptr {
	return _IOC(_IOC_READ|_IOC_WRITE, t, nr, size)
}

func _IOC_DIR(nr int) uintptr {
	return uintptr(((nr) >> _IOC_DIRSHIFT) & _IOC_DIRMASK)
}

func _IOC_TYPE(nr int) uintptr {
	return uintptr(((nr) >> _IOC_TYPESHIFT) & _IOC_TYPEMASK)
}

func _IOC_NR(nr int) uintptr {
	return uintptr(((nr) >> _IOC_NRSHIFT) & _IOC_NRMASK)
}

func _IOC_SIZE(nr int) uintptr {
	return uintptr(((nr) >> _IOC_SIZESHIFT) & _IOC_SIZEMASK)
}
