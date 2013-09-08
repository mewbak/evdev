// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

/*
evdev is a pure Go implementation of the Linux evdev API.
It allows a Go application to track events from any devices
mapped to `/dev/input/event[X]`.
*/
package evdev

/* References:

https://www.kernel.org/doc/Documentation/input/event-codes.txt
https://github.com/mirrors/linux-2.6/blob/f3b8436ad9a8ad36b3c9fa1fe030c7f38e5d3d0b/Documentation/input/ff.txt
*/
