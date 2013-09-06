## evdev

evdev is a pure Go implementation of the Linux evdev API.
It allows a Go application to track events from any devices
mapped to `/dev/input/event[X]`.


### Known issues

* Opening nodes in `/dev/input` requires root access. This means that
  our client applications do as well. This should be fixed.


### Usage

    go get github.com/jteeuwen/evdev


### License

Unless otherwise stated, all of the work in this project is subject to a
1-clause BSD license. Its contents can be found in the enclosed LICENSE file.

