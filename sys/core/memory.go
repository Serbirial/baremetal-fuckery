package core

// AllocateMemory assigns stack memory (and potentially heap/user memory)
func AllocateMemory(task *Task, stackSize uint32) {
	if stackSize == 0 {
		stackSize = DefaultStackSize
	}

	// Assign stack memory
	task.StackTop = AllocateStack(stackSize)
	task.StackBot = task.StackTop - uintptr(stackSize)

	// TODO: allocate heap or user memory if needed
}

// ProtectMemory sets up memory protection for per-task regions
func ProtectMemory(task *Task) {
	// Currently identity-mapped; no real protection yet
	// When implementing an MMU per-task, configure permissions here
}
