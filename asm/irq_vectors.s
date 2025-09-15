// asm/irq_vectors.s
// Minimal ARMv8 exception vector table
// Link with Go kernel

.section .text
.global vectors
vectors:
    // Reset / Undefined / SVC / Prefetch abort / Data abort / IRQ / FIQ
    b reset_handler           // Reset vector
    b default_handler         // Undefined instruction
    b default_handler         // SVC (Supervisor Call)
    b default_handler         // Prefetch abort
    b default_handler         // Data abort
    b irq_handler             // IRQ
    b default_handler         // FIQ

// Handlers
reset_handler:
    // Typically jump to Go kernel entry
    b kernel_main

default_handler:
    b default_handler  // Loop forever if unexpected exception

irq_handler:
    // Call Go IRQ handler
    bl go_irq_handler
    eret              // Return from exception
