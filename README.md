# Go Elgato StreamDeck interface (low level)

Code inspired from and partly taken from github.com/magicmonkey/go-streamdeck.

## Installation

On Linux, you might also need to add some `udev` rules.  Put this into `/etc/udev/rules.d/99-streamdeck.rules`:
```
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0060", MODE:="666"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0063", MODE:="666"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006c", MODE:="666"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="006d", MODE:="666"
```
