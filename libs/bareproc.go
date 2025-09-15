package libs

import "sys/core"

// CreateTask creates a new task with a function entry point and optional stack size
func CreateTask(entry func(), stackSize uint32) *core.Task {
	t := core.NewTask(entry, stackSize)
	core.AddTask(t) // add to scheduler queue
	return t
}

// RunTask marks a task ready so it can be picked up by any free core
func RunTask(t *core.Task) {
	if t.State != core.TaskFinished {
		t.State = core.TaskReady
		core.AddTask(t) // put it back in the ready queue if needed
	}
}
