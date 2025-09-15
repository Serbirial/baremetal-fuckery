// asm/wfi.s
// Simple wrapper to execute WFI instruction
// Call from Go: asmWFI()

.global asmWFI
asmWFI:
    wfi       // Wait For Interrupt
    ret
