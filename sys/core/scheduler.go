package core

import (
	"sync"
	"unsafe"
)

// MaxCores defines how many concurrent tasks we can run
const MaxCores = 8 // updated to support more cores

var (
	taskQueue       []*Task
	currentTask     [MaxCores]*Task
	coreReleased    [MaxCores]bool
	coreReleaseLock sync.Mutex
)

// InitScheduler initializes empty structures
func InitScheduler() {
	taskQueue = make([]*Task, 0, 16) // simple fixed-capacity queue
	for i := range currentTask {
		currentTask[i] = nil
		coreReleased[i] = false
	}
}

// ReleaseSecondaryCore signals a secondary core to start executing its scheduler
func ReleaseSecondaryCore(coreID int) {
	if coreID <= 0 || coreID >= MaxCores {
		return // core 0 is primary, secondary cores are 1+
	}
	coreReleaseLock.Lock()
	coreReleased[coreID] = true
	coreReleaseLock.Unlock()
}

// WaitForRelease blocks the core until it is released
func WaitForRelease(coreID int) {
	for {
		coreReleaseLock.Lock()
		released := coreReleased[coreID]
		coreReleaseLock.Unlock()
		if released {
			break
		}
		WFI() // wait for interrupt while spinning
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
