package core

import (
	"unsafe"
)

// MaxCores defines how many concurrent tasks we can run
const MaxCores = 4

var (
	taskQueue   []*Task
	currentTask [MaxCores]*Task
)

// InitScheduler initializes empty structures
func InitScheduler() {
	taskQueue = make([]*Task, 0, 16) // simple fixed-capacity queue
	for i := range currentTask {
		currentTask[i] = nil
	}
}

// AddTask adds a new task to the queue
func AddTask(t *Task) {
	taskQueue = append(taskQueue, t)
}

// Schedule returns the next task ready to run
func Schedule() *Task {
	for i, t := range taskQueue {
		if t.State == TaskReady || t.State == TaskRunning {
			// Remove from queue
			taskQueue = append(taskQueue[:i], taskQueue[i+1:]...)
			return t
		}
	}
	return nil
}

// SwitchToTask switches core to a new task
func SwitchToTask(coreID int, next *Task) {
	prev := currentTask[coreID]
	currentTask[coreID] = next
	if prev != nil {
		prev.State = TaskReady
	}
	next.State = TaskRunning
	asmSwitchTask(unsafe.Pointer(prev), unsafe.Pointer(next))
}

//go:linkname asmSwitchTask switch_task
func asmSwitchTask(prev, next unsafe.Pointer)
