package sys

import (
	"drivers"
	"libs/bareproc"
	"root"
	"sys/core"
	"sys/interrupts"
)

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

	// --- Create REPL as an isolated user-space task ---
	replTask := bareproc.CreateTask(root.Run, 0)
	bareproc.RunTask(replTask)

	// --- TODOS: start background kernel tasks ---
	// e.g., timers, device polling, etc.

	// --- Start per-core scheduling loops ---
	for i := 0; i < core.MaxCores; i++ {
		go core.RunCore(i) // each core handles its own tasks
	}

	// --- Kernel idle loop ---
	for {
		core.WFI() // idle until an interrupt occurs
	}
}
