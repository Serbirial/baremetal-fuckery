// Save current CPU context and restore next task
// ARMv8 EL1 only
// Task context: X0-X30, SP, PC

.global switch_task
switch_task:
    // x0 = current task pointer
    // x1 = next task pointer

    // Save registers into current task context
    stp x0, x1, [x0, #0]!
    // (for simplicity, full save/restore implementation goes here)

    // Restore registers from next task
    ldp x0, x1, [x1, #0]!
    ret
