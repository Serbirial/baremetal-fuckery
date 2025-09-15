/* linker_rpi.ld
 *
 * Linker script for a 64-bit Raspberry Pi kernel image.
 * The GPU firmware will load the binary at physical address 0x80000
 * when the file is named kernel8.img (or when arm_64bit=1).
 *
 * This script:
 *  - sets the load address (. = 0x80000)
 *  - places the boot assembly in .text.boot (so _start lives there)
 *  - exports __bss_start / __bss_end for the boot stub to zero memory
 *  - reserves a small stack area (you can grow it later)
 *
 * Keep this file in boot/ and reference from TinyGo or your Makefile.
 */

OUTPUT_FORMAT("elf64-littleaarch64")   /* 64-bit little-endian AArch64 ELF */
ENTRY(_start)                          /* program entry symbol (from boot.S) */

MEMORY
{
    /* We don't need a full memory map for linking a flat kernel,
     * but you can declare regions if you want to enforce constraints.
     * For now we place sections starting at the firmware load address.
     */
}

/* Start at the physical address the Pi firmware loads kernel8.img to */
SECTIONS
{
    /* Place the beginning of the image at 0x80000 (conventional for 64-bit kernels) */
    . = 0x80000;

    /* Keep the boot entry section first so _start is guaranteed at lowest text addr */
    .text.boot ALIGN(0x1000) : {
        KEEP(*(.text.boot))      /* our boot.S entrypoint goes here */
    }

    /* Main code/text */
    .text ALIGN(0x1000) : {
        *(.text*)                /* all other code (TinyGo/Go) */
    }

    /* Read-only data */
    .rodata ALIGN(0x1000) : {
        *(.rodata*)
    }

    /* Initialized data */
    .data ALIGN(0x1000) : {
        *(.data*)
    }

    /* Align before BSS and record its start */
    . = ALIGN(0x10);
    __bss_start = .;            /* symbol: start of bss (zero-init region) */

    /* Uninitialized data (BSS) - NOLOAD so no file bytes are produced for it */
    .bss (NOLOAD) : {
        *(.bss*)
        *(COMMON)
    }

    /* mark end of bss */
    __bss_end = .;

    /* Reserve a small stack area at the end of the linked output if you want to
     * reference addresses relative to the image end. Typically the boot stub will
     * set SP to some label (e.g. stack_top) that points into a pre-reserved .bss area.
     */

    /* Discard build metadata not needed on bare metal */
    /DISCARD/ : { *(.comment) *(.note*) *(.eh_frame*) }
}

/* Provide a human-friendly symbol of the image size (optional) */
PROVIDE(__kernel_image_size = ABSOLUTE(. - 0x80000));
