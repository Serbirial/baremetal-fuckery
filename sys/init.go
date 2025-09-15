package sys

import (
	"drivers"
	"libs/bareproc"
	"root"
	"sys/core"
	"sys/interrupts"
)

// MaxCores is the maximum cores your scheduler will handle
const MaxCores = 8

func KernelMain() {
	// --- Initialize low-level kernel systems ---

	// Initialize MMU (identity-mapped for now)
	InitMMU()

	// Initialize framebuffer
	FbInit()

	// Initialize system console (used by REPL and drivers)
	Con.Init()

	// Initialize interrupts (vector table, enable IRQs)
	interrupts.InitInterrupts()

	// --- Initialize all kernel-space drivers ---
	drivers.Keyboard.Init() // USB keyboard
	// TODO: add other drivers like timers, disk, etc. here

	// Print confirmation message
	Con.PutString("Kernel init complete: MMU + framebuffer + interrupts + USB OK\n", 0x00FF00FF)

	// --- Initialize the multi-core task scheduler ---
	core.InitScheduler()

	// --- Start secondary cores first ---
	for coreID := 1; coreID < MaxCores; coreID++ {
		// Release core by setting its entry point to RunCore(coreID)
		core.ReleaseSecondaryCore(coreID)
	}

	// --- Create REPL as an isolated user-space task on core 0 ---
	replTask := bareproc.CreateTask(root.Run, 0)
	bareproc.RunTask(replTask)

	// --- Run primary core scheduler loop on core 0 ---
	core.RunCore(0)

	// This point is never reached
	for {
		core.WFI()
	}
}
