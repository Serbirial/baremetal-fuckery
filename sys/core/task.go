package core

const (
	MaxTasks         = 16
	DefaultStackSize = 0x2000 // 8 KB
)

type TaskState int

const (
	TaskReady TaskState = iota
	TaskRunning
	TaskWaiting
	TaskFinished
)

type Task struct {
	ID       uint32
	Core     uint32    // Core currently running the task (if any)
	State    TaskState // Current task state
	StackTop uintptr   // Top of stack for context switching
	StackBot uintptr   // Base of stack (optional, for freeing)
	Entry    func()    // Task entry function
}

// Create a new Task with a stack
func NewTask(entry func(), stackSize uint32) *Task {
	if stackSize == 0 {
		stackSize = DefaultStackSize
	}

	stackTop := AllocateStack(stackSize)

	t := &Task{
		Entry:    entry,
		State:    TaskReady,
		StackTop: stackTop,
		StackBot: stackTop - uintptr(stackSize),
		Core:     0xFFFFFFFF, // Not yet assigned to any core
	}

	return t
}

// AllocateStack is a placeholder; replace with your actual memory allocator
func AllocateStack(size uint32) uintptr {
	// For now, just return a fake address
	return 0x80000000 // TODO: implement real memory allocation
}
