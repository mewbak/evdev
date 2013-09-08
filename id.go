// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

const (
	BusPCI       = 0x01
	BusISAPNP    = 0x02
	BusUSB       = 0x03
	BusHIL       = 0x04
	BusBluetooth = 0x05
	BusVirtual   = 0x06
	BusISA       = 0x10
	BusI8042     = 0x11
	BusXTKBD     = 0x12
	BusRS232     = 0x13
	BusGamePort  = 0x14
	BusParPort   = 0x15
	BusAmiga     = 0x16
	BusADB       = 0x17
	BusI2C       = 0x18
	BusHost      = 0x19
	BusGSC       = 0x1A
	BusAtari     = 0x1B
	BusSPI       = 0x1C
)

// Id represents the device identity.
//
// The bus type is the only field that contains accurate data.
// It can be compared to the BusXXX constants.
// The vendor, product and version fields are bus type-specific
// information relating to the identity of the device.
// Modern devices (typically using PCI or USB) do have information
// that can be used, but legacy devices (such as serial mice,
// PS/2 keyboards and game ports on ISA sound cards) do not.
// These numbers therefore are not meaningful for some
// values of bus type.
type Id struct {
	BusType uint16
	Vendor  uint16
	Product uint16
	Version uint16
}
