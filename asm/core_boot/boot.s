// boot.s
// ARMv8-A baremetal primary boot for multiple cores
// Link this first, sets up stack for core 0 and jumps to Go kernel

.global _start
_start:
    // Disable interrupts temporarily
    msr daifset, #0xF

    // Set up stack for core 0
    ldr x0, =stack0_top
    mov sp, x0

    // Call Go KernelMain
    bl KernelMain

    // If KernelMain ever returns, halt
halt:
    wfi
    b halt

// Reserve stack memory (8 KB per core)
.section .bss
.align 16
stack0:
    .space 8192
stack0_top:
