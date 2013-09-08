## evdev

**Note**: This is work in progress. Use at your own risk.

evdev is a pure Go implementation of the Linux evdev API.
It allows a Go application to track events from any devices
mapped to `/dev/input/event[X]`.


### Known issues

### Force Feedback

Not all of the force feedback API is implemented.
The documentation on this is pretty flaky, so it may take a while.


#### Permissions

Opening nodes in `/dev/input` may require root access. This means that
our client applications do as well. To solve this, there are a couple
of options.

The most sensible one is to use a `udev` rule to give device access
to anyone in the `input` group. Then add yourself to this group.
This hinges on the question whether or not your system uses `udev`.
For Arch Linux, `udev` comes pre-installed as a part of `systemd`.

Here is a short listing of the steps to undertake to make this work,
but we strongly advise that you read through the appropriate
[documentation](http://www.reactivated.net/writing_udev_rules.html)
on what `udev` rules are and how to safely create or edit them.

As root, perform the following steps:

	$ mkdir -p /etc/udev/rules.d
	$ nano /etc/udev/rules.d/99-input.rules

Put this in the file:

	KERNEL=="event*", NAME="input/%k", MODE="660", GROUP="input"

Save and exit nano. Then create the `input` group and add yourself to it:

	$ groupadd -f input
	$ gpasswd -a <YOURUSERNAME> input

This will add any input devices to the `input` group. Only users who are in
this group, will be able to read from them. Reboot your machine to make
these changes take effect.

`/dev/input' should now list someting like this:

	$ ls -l /dev/input/
	total 0
	drwxr-xr-x 2 root root     120 Sep  7 18:10 by-id
	drwxr-xr-x 2 root root     140 Sep  7 18:10 by-path
	crw-rw---- 1 root input 13, 64 Sep  7 18:10 event0
	crw-rw---- 1 root input 13, 65 Sep  7 18:10 event1
	crw-rw---- 1 root input 13, 74 Sep  7 18:10 event10
	crw-rw---- 1 root input 13, 75 Sep  7 18:10 event11
	crw-rw---- 1 root input 13, 76 Sep  7 18:10 event12
	crw-rw---- 1 root input 13, 77 Sep  7 18:10 event13
	crw-rw---- 1 root input 13, 78 Sep  7 18:10 event14
	crw-rw---- 1 root input 13, 66 Sep  7 18:10 event2
	crw-rw---- 1 root input 13, 67 Sep  7 18:10 event3
	crw-rw---- 1 root input 13, 68 Sep  7 18:10 event4
	crw-rw---- 1 root input 13, 69 Sep  7 18:10 event5
	crw-rw---- 1 root input 13, 70 Sep  7 18:10 event6
	crw-rw---- 1 root input 13, 71 Sep  7 18:10 event7
	crw-rw---- 1 root input 13, 72 Sep  7 18:10 event8
	crw-rw---- 1 root input 13, 73 Sep  7 18:10 event9
	crw-r----- 1 root root  13, 63 Sep  7 18:10 mice
	crw-r----- 1 root root  13, 32 Sep  7 18:10 mouse0


### Usage

    go get github.com/jteeuwen/evdev


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

