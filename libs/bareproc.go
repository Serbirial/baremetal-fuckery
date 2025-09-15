package libs

import "sys/core"

// CreateTask creates a new task with a function entry point and optional stack size.
// The task is added to the scheduler queue immediately.
func CreateTask(entry func(), stackSize uint32) *core.Task {
	t := core.NewTask(entry, stackSize)
	core.AddTask(t) // add to scheduler queue
	return t
}

// RunTask marks a task as ready so it can be picked up by any free core.
// If the task has finished, it will not be re-added.
func RunTask(t *core.Task) {
	if t.State != core.TaskFinished {
		t.State = core.TaskReady
		core.AddTask(t) // put it back in the ready queue
	}
}
