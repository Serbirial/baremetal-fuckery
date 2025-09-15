package sys

import (
	"drivers"
	"root"
)

func KernelMain() {
	// Initialize MMU (identity-mapped for now)
	InitMMU()

	// Initialize framebuffer
	FbInit()

	// Initialize system console (used by REPL and drivers)
	Con.Init()

	// Initialize interrupts (USB, timer, etc.)
	InitInterrupts()

	// Initialize USB keyboard driver
	drivers.Keyboard.Init()

	// Print confirmation message
	Con.PutString("MMU + framebuffer + interrupts + USB OK\n", 0x00FF00FF)

	// Start the REPL (interactive console)
	root.Run() // REPL now polls the real USB keyboard

	// Fallback infinite loop (REPL should not exit)
	for {
	}
}
