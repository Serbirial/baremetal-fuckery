// sys/mmu.go
package sys

import "unsafe"

// ARMv8 MMU registers (EL1)
const (
	SCTLR_MMU_EN = 1 << 0 // bit 0 enables MMU
)

// Page table aligned to 16KB
var pageTable [4096]uint64 // Level 1 translation table

// InitMMU sets up a basic identity-mapped MMU
func InitMMU() {
	setupPageTables()
	enableMMU()
}

// setupPageTables fills L1 table with 1:1 physical mappings
func setupPageTables() {
	// Map first 4GB of memory 1:1, 2MB sections
	for i := 0; i < 4096; i++ {
		pageTable[i] = (uint64(i) << 21) | 0x3 | (1 << 10) // AF=1, SH=Inner, AP=RW, AttrIdx=0
	}
}

// enableMMU writes TTBR0_EL1, TCR_EL1, SCTLR_EL1
func enableMMU() {
	ttbr := uintptr(unsafe.Pointer(&pageTable))
	setTTBR0(ttbr)
	setTCR()
	barrier()
	setSCTLR(SCTLR_MMU_EN)
}

// --- assembly helpers ---
func setTTBR0(addr uintptr)
func setTCR()
func barrier()
func setSCTLR(val uint64)
