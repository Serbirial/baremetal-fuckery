// secondary.s
// ARMv8-A secondary cores entry
// Each core starts here after being released by primary core

.global secondary_entry
secondary_entry:
    // x0 = core ID passed from primary (optional)
    mov x1, sp

    // Set stack for this core (8 KB per core)
    // Simple scheme: stackN = stack0 + 8192 * coreID
    ldr x2, =stack0_top
    mov x3, x0          // coreID
    lsl x3, x3, #13     // 8192 bytes per core
    add sp, x2, x3

    // Call Go RunCore(coreID)
    bl RunCore

    // Halt if RunCore returns (should never)
secondary_halt:
    wfi
    b secondary_halt

// Make sure stack memory covers multiple cores
.section .bss
.align 16
stack0:
    .space 8192 * 8   // reserve up to 8 cores
stack0_top:
