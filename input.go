// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

// Multitouch tools
const (
	MtToolFinger = 0
	MtToolPen    = 1
	MtToolMax    = 1
)

// Input device properties and quirks.
//
// Normally, userspace sets up an input device based on the data it emits,
// i.e., the event types. In the case of two devices emitting the same event
// types, additional information can be provided in the form of device
// properties.
const (
	/*	The InputPropPointer property indicates that the device is not transposed
		on the screen and thus requires use of an on-screen pointer to trace user's
		movements. Typical pointer devices: touchpads, tablets, mice; non-pointer
		device: touchscreen.

		If neither InputPropDirect or InputPropPointer are set, the property is
		considered undefined and the device type should be deduced in the
		traditional way, using emitted event types.
	*/
	InputPropPointer = 0x00

	/*	The InputPropDirect property indicates that device coordinates should be
		directly mapped to screen coordinates (not taking into account trivial
		transformations, such as scaling, flipping and rotating). Non-direct input
		devices require non-trivial transformation, such as absolute to relative
		transformation for touchpads. Typical direct input devices: touchscreens,
		drawing tablets; non-direct devices: touchpads, mice.

		If neither InputPropDirect or InputPropPointer are set, the property is
		considered undefined and the device type should be deduced in the
		traditional way, using emitted event types.
	*/
	InputPropDirect = 0x01

	/*	For touchpads where the button is placed beneath the surface, such that
		pressing down on the pad causes a button click, this property should be
		set. Common in clickpad notebooks and macbooks from 2009 and onwards.

		Originally, the buttonpad property was coded into the bcm5974 driver
		version field under the name integrated button. For backwards
		compatibility, both methods need to be checked in userspace.
	*/
	InputPropButtonPad = 0x02

	/*	Some touchpads, most common between 2008 and 2011, can detect the presence
		of multiple contacts without resolving the individual positions; only the
		number of contacts and a rectangular shape is known. For such
		touchpads, the semi-mt property should be set.

		Depending on the device, the rectangle may enclose all touches, like a
		bounding box, or just some of them, for instance the two most recent
		touches. The diversity makes the rectangle of limited use, but some
		gestures can normally be extracted from it.

		If InputPropSemiMT is not set, the device is assumed to be a true
		multi-touch device.
	*/
	InputPropSemiMT = 0x03

	InputPropMax   = 0x1f
	InputPropCount = InputPropMax + 1
)
