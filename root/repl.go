package root

import "drivers"

// Run starts the REPL loop
func Run() {
	// Initialize the REPL console
	Con.Init("> ")

	for {
		// Poll for keyboard input
		if key, ok := PollKey(); ok {
			Con.HandleKey(key)
		}
		// Optional: small delay to reduce busy-waiting
	}
}

// PollKey returns a single key press (ASCII) if available
func PollKey() (byte, bool) {
	if key, ok := drivers.Keyboard.Poll(); ok {
		return drivers.KeyToASCII(key, drivers.Keyboard.shift), true
	}
	return 0, false
}
