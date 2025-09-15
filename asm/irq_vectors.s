// asm/irq_vectors.s
// Minimal ARMv8 exception vector table
// Link with Go kernel

    .section .text
    .align 11               // 2048-byte alignment for VBAR_EL1
    .global irq_vector_table
irq_vector_table:
    b reset_handler           // Reset vector
    b default_handler         // Undefined instruction
    b default_handler         // SVC
    b default_handler         // Prefetch abort
    b default_handler         // Data abort
    b irq_handler             // IRQ
    b default_handler         // FIQ
    b default_handler         // Reserved / padding

// --- Handlers ---
reset_handler:
    // Jump to Go kernel
    bl KernelMain
    b .                       // halt if KernelMain returns

default_handler:
    b .                       // infinite loop

irq_handler:
    bl go_irq_handler
    subs pc, lr, #4           // return from IRQ
