// sys/interrupts.go
package sys

import "unsafe"

// InitInterrupts sets the vector table and enables IRQs
func InitInterrupts() {
	asmSetVBAR(uintptr(unsafe.Pointer(&irq_vector_table)))
	EnableInterrupts()
}

// EnableInterrupts clears IRQ mask
func EnableInterrupts() {
	asmEnableIRQ()
}

// Go IRQ handler called from assembly
//export go_irq_handler
func go_irq_handler() {
	// TODO: handle multiple devices, for now timer only
	handleTimerIRQ()
}

// Go synchronous exception handler
//export go_sync_handler
func go_sync_handler() {
	// Just halt for now
	for {
	}
}

// Example timer IRQ handler
func handleTimerIRQ() {
	// Future: clear timer, update counters, etc.
}

// Link Go names to the ASM implementations
//go:linkname asmEnableIRQ asmEnableIRQ
func asmEnableIRQ()

//go:linkname asmSetVBAR asmSetVBAR
func asmSetVBAR(addr uintptr)

// Declare irq_vector_table from assembly
//go:linkname irq_vector_table irq_vector_table
var irq_vector_table [0]byte
