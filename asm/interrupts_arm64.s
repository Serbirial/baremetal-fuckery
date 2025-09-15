// asm/interrupts_arm64.s
// ARMv8 EL1 vector table and basic IRQ setup
// Universal for ARMv8, calls Go handlers

.global irq_vector_table
irq_vector_table:
    // Each vector must be 128-byte aligned (ARMv8 requirement)
    b el1_sync        // 0x000: Synchronous EL1
    b el1_irq         // 0x080: IRQ EL1
    b el1_fiq         // 0x100: FIQ EL1
    b el1_serror      // 0x180: SError EL1

.align 7
el1_sync:
    // Save registers
    stp x0, x1, [sp, #-16]! 
    bl go_sync_handler
    ldp x0, x1, [sp], #16
    eret

el1_irq:
    stp x0, x1, [sp, #-16]!
    bl go_irq_handler
    ldp x0, x1, [sp], #16
    eret

el1_fiq:
    // Not implemented, just return
    eret

el1_serror:
    // Not implemented, just return
    eret

// Enable IRQs
.global asmEnableIRQ
asmEnableIRQ:
    msr daifclr, #2   // Clear I-bit to enable IRQ
    ret

// Set vector base address
.global asmSetVBAR
asmSetVBAR:
    msr vbar_el1, x0  // Set EL1 vector base
    dsb sy
    isb
    ret
