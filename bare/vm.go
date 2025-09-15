// bare/vm.go
package bare

import "errors"

// VM represents the virtual machine state
type VM struct {
	IP    int     // instruction pointer
	Code  []byte  // bytecode
	Stack []int32 // simple integer stack
}

// Max stack size
const StackSize = 256

// NewVM creates a new VM instance
func NewVM() *VM {
	return &VM{
		Stack: make([]int32, 0, StackSize),
	}
}

// Run executes the loaded bytecode
func (vm *VM) Run() error {
	for vm.IP < len(vm.Code) {
		op := vm.Code[vm.IP]
		switch op {
		case 0x01: // PUSH
			if vm.IP+4 >= len(vm.Code) {
				return errors.New("unexpected end of bytecode")
			}
			val := int32(vm.Code[vm.IP+1]) | int32(vm.Code[vm.IP+2])<<8 | int32(vm.Code[vm.IP+3])<<16 | int32(vm.Code[vm.IP+4])<<24
			vm.Stack = append(vm.Stack, val)
			vm.IP += 5
		case 0x02: // ADD
			if len(vm.Stack) < 2 {
				return errors.New("stack underflow")
			}
			a := vm.Stack[len(vm.Stack)-1]
			b := vm.Stack[len(vm.Stack)-2]
			vm.Stack = vm.Stack[:len(vm.Stack)-2]
			vm.Stack = append(vm.Stack, a+b)
			vm.IP++
		case 0x03: // SUB
			if len(vm.Stack) < 2 {
				return errors.New("stack underflow")
			}
			a := vm.Stack[len(vm.Stack)-1]
			b := vm.Stack[len(vm.Stack)-2]
			vm.Stack = vm.Stack[:len(vm.Stack)-2]
			vm.Stack = append(vm.Stack, b-a)
			vm.IP++
		case 0x04: // PRINT
			if len(vm.Stack) < 1 {
				return errors.New("stack underflow")
			}
			val := vm.Stack[len(vm.Stack)-1]
			vm.Stack = vm.Stack[:len(vm.Stack)-1]
			PrintInt(val) // use OS console function
			vm.IP++
		default:
			return errors.New("unknown opcode")
		}
	}
	return nil
}
