// asm/mmu_arm64.s
// Minimal ARMv8 EL1 MMU enable helpers for TinyGo

// ---------------------------
// Set TTBR0_EL1 (Translation Table Base Register 0)
// ---------------------------
.global setTTBR0          // make symbol visible to Go
setTTBR0:
    // x0 contains the physical address of our Level 1 page table
    msr ttbr0_el1, x0    // write TTBR0_EL1 = x0
    ret                  // return to caller

// ---------------------------
// Set TCR_EL1 (Translation Control Register)
// ---------------------------
.global setTCR
setTCR:
    // TCR_EL1 controls:
    // - virtual address size
    // - granule size (4KB here)
    // - caching / shareability attributes
    // Value 0xD0 is minimal: 48-bit VA, 4KB granule, inner shareable
    mov x0, #0x00000000000000D0
    msr tcr_el1, x0       // write TCR_EL1 = x0
    ret                   // return

// ---------------------------
// Memory barriers
// ---------------------------
.global barrier
barrier:
    dsb sy                // Data Synchronization Barrier (ensure memory ops complete)
    isb                   // Instruction Synchronization Barrier (flush pipeline)
    ret

// ---------------------------
// Set SCTLR_EL1 (System Control Register) to enable MMU
// ---------------------------
.global setSCTLR
setSCTLR:
    // x0 contains the SCTLR flags
    // We will set bit 0 = MMU enable (SCTLR_MMU_EN)
    msr sctlr_el1, x0     // write SCTLR_EL1 = x0
    isb                    // flush pipeline to ensure new settings take effect immediately
    ret
