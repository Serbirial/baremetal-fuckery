// root/repl.go
package root

import "drivers"

//go:linkname asmWFI asmWFI
func asmWFI()

// Run starts the REPL loop
func Run() {
	// Initialize the REPL console with a prompt
	Con.Init("> ")

	// Initialize the USB keyboard
	drivers.Keyboard.Init()

	for {
		// Poll the USB keyboard
		if key, ok := PollKey(); ok {
			Con.HandleKey(key)
		}

		// Wait for the next interrupt to reduce CPU usage
		asmWFI()
	}
}

// PollKey returns a single ASCII key press from the USB keyboard
func PollKey() (byte, bool) {
	if hid, pressed := drivers.Keyboard.Poll(); pressed {
		return drivers.KeyToASCII(hid, drivers.Keyboard.shift), true
	}
	return 0, false
}
