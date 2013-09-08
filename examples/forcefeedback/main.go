// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/jteeuwen/evdev"
	"os"
	"os/signal"
)

func main() {
	node := parseArgs()

	// Create and open our device.
	dev, err := evdev.Open(node)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	// Make sure it is closed once we are done.
	defer dev.Close()

	// Ensure this device supports the needed event types.
	if !dev.Test(dev.EventTypes(), evdev.EvForceFeedback, evdev.EvForceFeedbackStatus) {
		fmt.Fprintf(os.Stderr, "Device %q does not support force feedback events.\n", node)
		return
	}

	// Set effect gain factor, to ensure the effect strength is
	// the same on all FF devices we may be working with.
	dev.SetEffectGain(75) // 75%

	// List Force Feedback capabilities
	listCapabilities(dev)

	// Create, upload and play some effects.
	setEffects(dev)

	// Wait for incoming events or exit signals.
	poll(dev)
}

// poll waits for incoming events or exit signals.
//
// Events are triggered whenever an effect's state is altered.
// This only applies to devices with EvForceFeedbackStatus support.
func poll(dev *evdev.Device) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	for {
		select {
		case <-signals:
			return

		case evt := <-dev.Inbox:
			if evt.Type != evdev.EvForceFeedbackStatus {
				continue
			}

			if evt.Value == evdev.FFStatusStopped {
				fmt.Printf("%v effect %d is now stopped\n", evt.Time, evt.Code)
			} else {
				fmt.Printf("%v effect %d is now playing\n", evt.Time, evt.Code)
			}
		}
	}
}

// setEffects creates, uploads and plays a new Force feedback effect.
// This function uploads only 1 effect, but it can deal with
// up to N effects at the same time. Where N is whatever value
// returned from `Device.ForceFeedbackCaps()`.
func setEffects(dev *evdev.Device) {
	_, caps := dev.ForceFeedbackCaps()

	var effect evdev.Effect
	effect.Id = -1
	effect.Trigger.Button = 0
	effect.Trigger.Interval = 0
	effect.Replay.Length = 20000 // 20 seconds
	effect.Replay.Delay = 0

	// Some samples of various effect types.
	// (un)comment any one to try them out. Note that
	// the device must support a given effect type.

	switch {
	case dev.Test(caps, evdev.FFRumble):
		rumble(&effect)
	case dev.Test(caps, evdev.FFPeriodic):
		periodic(&effect)
	case dev.Test(caps, evdev.FFConstant):
		constant(&effect)
	case dev.Test(caps, evdev.FFSpring):
		spring(&effect)
	case dev.Test(caps, evdev.FFDamper):
		damper(&effect)
	}

	// Upload the effect.
	dev.SetEffects(&effect)

	fmt.Printf("Effect id: %d\n", effect.Id)

	// Play the effect.
	dev.PlayEffect(effect.Id)

	// Delete the effect.
	dev.UnsetEffects(&effect)
}

func damper(e *evdev.Effect) {
	e.Type = evdev.FFDamper
	e.Direction = 0x6000 // 135 degrees

	var c [2]evdev.ConditionEffect
	c[0].RightSaturation = 0x7fff
	c[0].LeftSaturation = 0x7fff
	c[0].RightCoeff = 0x2000
	c[0].LeftCoeff = 0x2000
	c[0].Deadband = 0x0
	c[0].Center = 0x0
	c[1] = c[0]
	e.SetData(c)
}

func spring(e *evdev.Effect) {
	e.Type = evdev.FFSpring
	e.Direction = 0x6000 // 135 degrees

	var c [2]evdev.ConditionEffect
	c[0].RightSaturation = 0x7fff
	c[0].LeftSaturation = 0x7fff
	c[0].RightCoeff = 0x2000
	c[0].LeftCoeff = 0x2000
	c[0].Deadband = 0x0
	c[0].Center = 0x0
	c[1] = c[0]
	e.SetData(c)
}

func constant(e *evdev.Effect) {
	e.Type = evdev.FFConstant
	e.Direction = 0x6000 // 135 degrees

	var c evdev.ConstantEffect
	c.Level = 0x2000 // Strength : 25 %
	c.Envelope.AttackLength = 0x100
	c.Envelope.AttackLevel = 0
	c.Envelope.FadeLength = 0x100
	c.Envelope.FadeLevel = 0
	e.SetData(c)
}

func rumble(e *evdev.Effect) {
	e.Type = evdev.FFRumble

	var r evdev.RumbleEffect
	r.StrongMagnitude = 0
	r.WeakMagnitude = 0xc000
	e.SetData(r)
}

func periodic(e *evdev.Effect) {
	e.Type = evdev.FFPeriodic
	e.Direction = evdev.DirLeft // Along X axis

	var p evdev.PeriodicEffect
	p.Waveform = evdev.FFSine
	p.Period = 26        // 0.1*0x100 = 0.1 second
	p.Magnitude = 0x4000 // 0.5 * Maximum magnitude
	p.Offset = 0
	p.Phase = 0
	p.Envelope.AttackLength = 0x100
	p.Envelope.AttackLevel = 0
	p.Envelope.FadeLength = 0x100
	p.Envelope.FadeLevel = 0

	e.SetData(p)
}

// listCapabilities lists Force feedback capabilities for a given device.
//
// Testing for individual effect types can be done using the
// Device.Supports() method.
func listCapabilities(dev *evdev.Device) {
	// Fetch the force feedback capabilities.
	// The number of simultaneous effects and a
	// bitset describing the type of effects.
	count, caps := dev.ForceFeedbackCaps()

	fmt.Printf("Number of simultaneous effects: %d\n", count)

	for n := 0; n < caps.Len(); n++ {
		if !caps.Test(n) {
			continue
		}

		fmt.Printf(" - Effect 0x%02x: ", n)

		switch n {
		case evdev.FFConstant:
			fmt.Printf("Constant")
		case evdev.FFPeriodic:
			fmt.Printf("Periodic")
		case evdev.FFSpring:
			fmt.Printf("Spring")
		case evdev.FFFriction:
			fmt.Printf("Friction")
		case evdev.FFRumble:
			fmt.Printf("Rumble")
		case evdev.FFDamper:
			fmt.Printf("Damper")
		case evdev.FFRamp:
			fmt.Printf("Ramp")
		}

		fmt.Println()
	}
}

func parseArgs() string {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <node>\n", os.Args[0])
		os.Exit(1)
	}

	return flag.Args()[0]
}
