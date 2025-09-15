/* boot/boot.S
 * This is the very first code the CPU runs after the Pi's firmware
 * loads our kernel image at address 0x80000.
 * Its job:
 *   1. Set up the stack pointer (SP)
 *   2. Clear the BSS section (zero-initialized data)
 *   3. Jump into our Go entrypoint: kernel_main()
 */

.section .text.boot     // Put this code into a special boot section
.global _start          // Make the symbol "_start" visible to the linker

_start:
    // ==============================
    // 1. Initialize the stack pointer
    // ==============================

    // Load the address of "stack_top" into register x0
    // This is the top of our reserved stack space.
    ldr x0, =stack_top

    // Move that value into the SP (stack pointer) register
    // Now any push/pop or local variables will have a valid stack.
    mov sp, x0

    // ==============================
    // 2. Clear the BSS section
    // ==============================

    // BSS is memory for variables that should start as 0.
    // The linker gives us two symbols: __bss_start and __bss_end.

    // Load start address of BSS into x1
    ldr x1, =__bss_start

    // Load end address of BSS into x2
    ldr x2, =__bss_end

1:  // Loop label: zero memory until we reach the end
    cmp x1, x2           // Compare current pointer (x1) with end (x2)
    b.hs 2f              // If x1 >= x2, we’re done → jump to label 2

    str xzr, [x1], #8    // Store zero (xzr is always 0) at [x1]
                         // Then post-increment x1 by 8 bytes (ARM64 is 64-bit)

    b 1b                 // Repeat the loop

2:
    // ==============================
    // 3. Call into Go kernel_main()
    // ==============================

    // Branch with link (bl) to kernel_main()
    // This jumps into our Go code, but also stores return address in LR.
    bl kernel_main

    // ==============================
    // 4. Fallback: hang forever
    // ==============================

hang:
    wfe                  // Wait For Event (low-power idle)
    b hang               // Infinite loop
                         // If kernel_main ever returns, we stay here.


// ==============================
// Reserve stack memory
// ==============================

// Place a chunk of zero-initialized memory for stack.
// The stack grows downward, so we use the *end* as the start.

.section .bss
    .space 4096          // Reserve 4 KB for stack
stack_top:
