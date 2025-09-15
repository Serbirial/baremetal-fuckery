// sys/interrupts/interrupts.go
package interrupts

var coreReleaseFlags [16]uint32 // support up to 16 cores for now

// InitInterrupts initializes the vector table and enables IRQs
func InitInterrupts() {
	// TODO: set vector table base address, enable IRQs
	for i := range coreReleaseFlags {
		coreReleaseFlags[i] = 0
	}
}

// HandleIRQ is called from the assembly IRQ vector
func HandleIRQ() {
	// TODO: check source of IRQ and dispatch
}

// ReleaseCore signals a secondary core to start its scheduler
func ReleaseCore(coreID int) {
	if coreID <= 0 || coreID >= len(coreReleaseFlags) {
		return
	}
	coreReleaseFlags[coreID] = 1
}

// PollReleaseFlag is polled by secondary cores during early boot
func PollReleaseFlag(coreID int) bool {
	if coreID <= 0 || coreID >= len(coreReleaseFlags) {
		return false
	}
	if coreReleaseFlags[coreID] != 0 {
		coreReleaseFlags[coreID] = 0 // clear after release
		return true
	}
	return false
}
