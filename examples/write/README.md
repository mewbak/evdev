## Write

This program demonstrates how to write to an `evdev` device.
It loads up a keyboard and toggles the LED states on it:

* Turn on CAPS_LOCK for 200 ms then off again.
* Turn on NUM_LOCK for 200 ms then off again.
* Turn on SCROLL_LOCK for 200 ms then off again.
* Repeat from the top.

### Usage

	$ go build
	$ ./write /dev/input/event0

