// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import "unsafe"

// Values describing the status of a force-feedback effect
const (
	FFStatusStopped = 0x00
	FFStatusPlaying = 0x01
	FFStatusMax     = 0x01
)

// Force feedback effect types.
const (
	FFRumble    = 0x50
	FFPeriodic  = 0x51
	FFConstant  = 0x52
	FFSpring    = 0x53
	FFFriction  = 0x54
	FFDamper    = 0x55
	FFInertia   = 0x56
	FFRamp      = 0x57
	FFEffectMin = FFRumble
	FFEffectMax = FFRamp
)

// Force feedback periodic effect types
const (
	FFSquare      = 0x58
	FFTriangle    = 0x59
	FFSine        = 0x5a
	FFSawUp       = 0x5b
	FFSawDown     = 0x5c
	FFCustom      = 0x5d
	FFWaveformMin = FFSquare
	FFWaveformMax = FFCustom
)

// Set Force feedback device properties
const (
	FFGain       = 0x60
	FFAutoCenter = 0x61
	FFMax        = 0x7f
	FFCount      = FFMax + 1
)

// Directions encoded in Effect.Direction
const (
	DirDown  = 0x0000 // 0 degrees
	DirLeft  = 0x4000 // 90 degrees
	DirUp    = 0x8000 // 180 degrees
	DirRight = 0xc000 // 270 degrees
)

/*
Effect describes any of the supported Force Feedback effects.

Supported effects are as follows:

	- FFConstant        Renders constant force effects
	- FFPeriodic        Renders periodic effects with the following waveforms:
	  - FFSquare        Square waveform
	  - FFTriangle      Triangle waveform
	  - FFSine          Sine waveform
	  - FFSawUp         Sawtooth up waveform
	  - FFSawDown       Sawtooth down waveform
	  - FFCustom        Custom waveform
	- FFRamp            Renders ramp effects
	- FFSpring          Simulates the presence of a spring
	- FFFriction        Simulates friction
	- FFDamper          Simulates damper effects
	- FFRumble          Rumble effects
	- FFInertia         Simulates inertia
	- FFGain            Gain is adjustable
	- FFAutoCenter      Autocenter is adjustable

Note: In most cases you should use FFPeriodic instead of FFRumble. All
devices that support FFRumble support FFPeriodic (square, triangle,
sine) and the other way around.

Note: The exact layout of FFCustom waveforms is undefined for the
time being as no driver supports it yet.

Note: All duration values are expressed in milliseconds.
Values above 32767 ms (0x7fff) should not be used and have unspecified results.
*/
type Effect struct {
	Type      uint16
	Id        int16
	Direction uint16
	Trigger   Trigger
	Replay    Replay
	data      unsafe.Pointer
}

// Data returns the event data structure as a concrete type.
// Its type depends on the value of Effect.Type and can be any of:
//
//    FFConstant -> ConstantEffect
//    FFPeriodic -> PeriodicEffect
//    FFRamp     -> RampEffect
//    FFRumble   -> RumbleEffect
//    FFSpring   -> [2]ConditionEffect
//    FFDamper   -> [2]ConditionEffect
//
// This returns nil if the type was not recognized.
func (e *Effect) Data() interface{} {
	// FIXME(jimt): Deal with: FFFriction, FFInertia:
	// Unsure what they should return.

	if e.data == nil {
		return nil
	}

	switch e.Type {
	case FFConstant:
		return *(*ConstantEffect)(e.data)
	case FFPeriodic:
		return *(*PeriodicEffect)(e.data)
	case FFRamp:
		return *(*RampEffect)(e.data)
	case FFRumble:
		return *(*RumbleEffect)(e.data)
	case FFSpring, FFDamper:
		return *(*[2]ConditionEffect)(e.data)
	}

	return nil
}

// SetData sets the event data structure.
func (e *Effect) SetData(v interface{}) {
	if v != nil {
		e.data = unsafe.Pointer(&v)
	}
}

type Replay struct {
	Length uint16
	Delay  uint16
}

type Trigger struct {
	Button   uint16
	Interval uint16
}

type Envelope struct {
	AttackLength uint16
	AttackLevel  uint16
	FadeLength   uint16
	FadeLevel    uint16
}

// ConstantEffect renders constant force-feedback effects.
type ConstantEffect struct {
	Level    int16
	Envelope Envelope
}

// RampEffect renders ramp force-feedback effects.
type RampEffect struct {
	StartLevel int16
	EndLevel   int16
	Envelope   Envelope
}

// ConditionEffect represents a confitional force feedback effect.
type ConditionEffect struct {
	RightSaturation uint16
	LeftSaturation  uint16
	RightCoeff      int16
	LeftCoeff       int16
	Deadband        uint16
	Center          int16
}

// PeriodicEffect renders periodic force-feedback effects with
// the following waveforms: Square, Triangle, Sine, Sawtooth
// or a custom waveform.
type PeriodicEffect struct {
	Waveform  uint16
	Period    uint16
	Magnitude int16
	Offset    int16
	Phase     uint16
	Envelope  Envelope

	custom_len  uint32
	custom_data unsafe.Pointer // *int16
}

// Data returns custom waveform information.
// This comes in the form of a signed 16-bit slice.
//
// The exact layout of a custom waveform is undefined for the
// time being as no driver supports it yet.
func (e *PeriodicEffect) Data() []int16 {
	if e.custom_data == nil {
		return nil
	}
	return (*(*[1<<31 - 1]int16)(e.custom_data))[:e.custom_len]
}

// SetData sets custom waveform information.
//
// The exact layout of a custom waveform is undefined for the
// time being as no driver supports it yet.
func (e *PeriodicEffect) SetData(v []int16) {
	e.custom_len = uint32(len(v))
	e.custom_data = unsafe.Pointer(nil)

	if e.custom_len > 0 {
		e.custom_data = unsafe.Pointer(&v[0])
	}
}

// The rumble effect is the most basic effect, it lets the
// device vibrate. The API contains support for two motors,
// a strong one and a weak one, which can be controlled independently.
type RumbleEffect struct {
	StrongMagnitude uint16
	WeakMagnitude   uint16
}

// ForceFeedbackCaps returns a bitset which specified the kind of Force Feedback
// effects supported by this device. The bits can be compared against
// the FFXXX constants. Additionally, it returns the number of effects
// this device can handle simultaneously.
//
// This is only applicable to devices with EvForceFeedback event support.
func (d *Device) ForceFeedbackCaps() (int, Bitset) {
	bs := NewBitset(24)
	buf := bs.Bytes()
	ioctl(d.fd.Fd(), _EVIOCGBIT(EvForceFeedback, len(buf)), unsafe.Pointer(&buf[0]))

	var count int32
	ioctl(d.fd.Fd(), _EVIOCGEFFECTS, unsafe.Pointer(&count))
	return int(count), bs
}

// SetEffects sends the given list of Force Feedback effects
// to the device. The length of the list should not exceed the
// count returned from `Device.ForceFeedbackCaps()`.
//
// After this call completes, the effect.Id field will contain
// the effect's id which must be used when playing or stopping the effect.
// It is also possible to reupload the same effect with the same
// id later on with new parameters. This allows us to update a
// running effect, without first stopping it.
//
// This is only applicable to devices with EvForceFeedback event support.
func (d *Device) SetEffects(list ...*Effect) bool {
	for _, effect := range list {
		err := ioctl(d.fd.Fd(), _EVIOCSFF, unsafe.Pointer(effect))
		if err != nil {
			return false
		}
	}

	return true
}

// UnsetEffects deletes the given effects from the device.
// This makes room for new effects in the device's memory.
// Note that this also stops the effect if it was playing.
//
// This is only applicable to devices with EvForceFeedback event support.
func (d *Device) UnsetEffects(list ...*Effect) bool {
	for _, effect := range list {
		err := ioctl(d.fd.Fd(), _EVIOCRMFF, int(effect.Id))
		if err != nil {
			return false
		}
	}

	return true
}

// SetEffectGain changes the force feedback gain.
//
// Not all devices have the same effect strength. Therefore,
// users should set a gain factor depending on how strong they
// want effects to be. This setting is persistent across
// access to the driver.
//
// The specified gain should be in the range 0-100.
// This is only applicable to devices with EvForceFeedback event support.
func (d *Device) SetEffectGain(gain int) {
	d.setEffectFactor(gain, FFGain)
}

// SetEffectAutoCenter changes the force feedback autocenter factor.
// The specified factor should be in the range 0-100.
// A value of 0 means: no autocenter.
//
// This is only applicable to devices with EvForceFeedback event support.
func (d *Device) SetEffectAutoCenter(factor int) {
	d.setEffectFactor(factor, FFAutoCenter)
}

// setEffectFactor changes the given effect factor.
// E.g.: FFAutoCenter or FFGain.
//
// This is only applicable to devices with EvForceFeedback event support.
func (d *Device) setEffectFactor(factor int, code uint16) {
	if factor < 0 {
		factor = 0
	}

	if factor > 100 {
		factor = 100
	}

	var e Event
	e.Type = EvForceFeedback
	e.Code = code
	e.Value = 0xffff * int32(factor) / 100
	d.Outbox <- e
}

// PlayEffect plays a previously uploaded effect.
func (d *Device) PlayEffect(id int16) {
	d.toggleEffect(id, true)
}

// StopEffect stops a previously uploaded effect from playing.
func (d *Device) StopEffect(id int16) {
	d.toggleEffect(id, false)
}

// ToggleEffect plays or stops a previously uploaded effect with the given id.
func (d *Device) toggleEffect(id int16, play bool) {
	var e Event
	e.Type = EvForceFeedback
	e.Code = uint16(id)

	if play {
		e.Value = 1
	} else {
		e.Value = 0
	}

	d.Outbox <- e
}
