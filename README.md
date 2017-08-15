# Ubuntu Rotater

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Travis](https://travis-ci.org/ramjac/ubuntu-rotater.svg?branch=master)](https://travis-ci.org/ramjac/ubuntu-rotater)
[![Go Report Card](https://goreportcard.com/badge/github.com/ramjac/ubuntu-rotater)](https://goreportcard.com/report/github.com/ramjac/ubuntu-rotater)

This is a simple program to rotate the display and input based on the output of the "monitor-sensors" command. It was written and tested on Ubuntu with a Dell Latitude 7275.

Currently it is hard coded to work with the "Wacom HID 4822" devices. If you've got a Wacom digitizer, you can check the device name with this command:

```
xsetwacom --list devices
```

If someone out there is interested in seeing this made into something that maybe take command line arguments so it can work with more devices, then send me a message and I'll se what I can do.

## Prerequisites

I'm using three commands that should work from the terminal: xrandr, xsetwacom, and monitor-sensors



