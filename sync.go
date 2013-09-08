// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

// Synchronization event values are undefined.
// Their usage is defined only by when they are
// sent in the evdev event stream.
//
// SynReport is used to synchronize and separate
// events into packets of input data changes occurring
// at the same moment in time. For example, motion
// of a mouse may set the RelX and RelY values for
// one motion, then emit a SynReport. The next motion
// will emit more RelX and RelY values and send
// another SynReport.
//
// SynConfig: to be determined.
//
// SynMTReport is used to synchronize and separate
// touch events. See the multi-touch-protocol.txt
// document for more information.
//
// SynDropped is used to indicate buffer overrun
// in the evdev client's event queue.
// Client should ignore all events up to and
// including next SynReport event and query the
// device (using EVIOCG* ioctls) to obtain its current state.
const (
	SynReport = iota
	SynConfig
	SynMTReport
	SynDropped
)
