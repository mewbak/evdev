// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

// Switch events
const (
	SwLid                = 0x00        // set = lid shut
	SwTabletMode         = 0x01        // set = tablet mode
	SwHeadphoneInsert    = 0x02        // set = inserted
	SwRFKillAll          = 0x03        // rfkill master switch, type "any"; set = radio enabled
	SwRadio              = SwRFKillAll // deprecated
	SwMicrophoneInsert   = 0x04        // set = inserted
	SwDock               = 0x05        // set = plugged into dock
	SwLineoutInsert      = 0x06        // set = inserted
	SwJackPhysicalInsert = 0x07        // set = mechanical switch set
	SwVideoOutInsert     = 0x08        // set = inserted
	SwCameraLensCover    = 0x09        // set = lens covered
	SwKeypadSlide        = 0x0a        // set = keypad slide out
	SwFrontProximity     = 0x0b        // set = front proximity sensor active
	SwRotateLock         = 0x0c        // set = rotate locked/disabled
	SwLineInInsert       = 0x0d        // set = inserted
	SwMax                = 0x0f
	SwCount              = SwMax + 1
)
