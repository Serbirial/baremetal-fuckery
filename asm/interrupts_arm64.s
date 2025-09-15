// asm/interrupts_arm64.s
// Minimal ARMv8 EL1 vector table with proper 128-byte spacing
// Universal for ARMv8, calls Go handlers

.global irq_vector_table
irq_vector_table:
    b el1_sync        // 0x000: Synchronous EL1
    .space 124        // pad to 128 bytes

    b el1_irq         // 0x080: IRQ EL1
    .space 124

    b el1_fiq         // 0x100: FIQ EL1
    .space 124

    b el1_serror      // 0x180: SError EL1
    .space 124

.align 7
// --- Vector handlers ---
el1_sync:
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

// --- Enable IRQs ---
.global asmEnableIRQ
asmEnableIRQ:
    msr daifclr, #2   // Clear I-bit to enable IRQ
    ret

// --- Set vector base address ---
.global asmSetVBAR
asmSetVBAR:
    msr vbar_el1, x0  // Set EL1 vector base
    dsb sy
    isb
    ret
