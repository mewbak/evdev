## Exclusive

This program demonstrates how to obtain an exclusive lock
on a given device. This ensures that you alone will receive
events from it. Other processes will not.

This example will attempt to obtain a device lock for 5 seconds.
Then releases it and exits.


### Usage

	$ go build
	$ ./exclusive /dev/input/event0

