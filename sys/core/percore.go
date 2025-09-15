package core

// RunCore runs a scheduling loop for a single CPU core
func RunCore(coreID int) {
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
