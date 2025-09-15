package core

import "sys/interrupts"

// RunCore runs the scheduling loop for a single CPU core.
// Core 0 (primary) releases secondary cores before entering its loop.
// Secondary cores spin in WFI until they are released.
func RunCore(coreID int) {
	if coreID != 0 {
		// Secondary cores wait until primary signals them
		for !interrupts.PollReleaseFlag(coreID) {
			WFI()
		}
	} else {
		// Primary core releases other cores
		ReleaseSecondaryCores()
	}

	// Start per-core scheduling loop
	for {
		t := Schedule()
		if t != nil {
			SwitchToTask(coreID, t)
		} else {
			// Idle until an interrupt occurs
			WFI()
		}
	}
}

// WFI calls the assembly wrapper for the ARMv8 WFI instruction
func WFI() {
	asmWFI()
}

//go:linkname asmWFI asmWFI
func asmWFI()

// ReleaseSecondaryCores signals all secondary cores to start scheduling
func ReleaseSecondaryCores() {
	for i := 1; i < MaxCores; i++ {
		interrupts.ReleaseCore(i)
	}
}
