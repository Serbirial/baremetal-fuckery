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

// nextStackAddr tracks the next free stack memory
var nextStackAddr uintptr = 0x80000000

// Create a new Task with stack allocation and memory setup
func NewTask(entry func(), stackSize uint32) *Task {
	if stackSize == 0 {
		stackSize = DefaultStackSize
	}

	t := &Task{
		Entry: entry,
		State: TaskReady,
		Core:  0xFFFFFFFF, // Not yet assigned
	}

	// Allocate memory for this task
	AllocateMemory(t, stackSize)
	ProtectMemory(t)

	return t
}

// AllocateStack returns the top of a new stack and updates nextStackAddr
func AllocateStack(size uint32) uintptr {
	addr := nextStackAddr
	nextStackAddr += uintptr(size)
	return addr + uintptr(size) // return top of stack
}
