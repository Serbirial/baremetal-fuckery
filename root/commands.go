// root/commands.go
package root

import (
	"bare"
	"sys"
)

// CommandFunc defines a function that executes a command
type CommandFunc func(args []string)

// Built-in commands map (declare first, fill later)
var commands map[string]CommandFunc

// init fills the commands map
func init() {
	commands = map[string]CommandFunc{
		"help":  CmdHelp,
		"print": CmdPrint,
		"add":   CmdAdd,
	}
}

// ExecuteCommand parses a line and runs the corresponding command
func ExecuteCommand(line string) {
	line = bare.Trim(line)
	if line == "" {
		return
	}

	parts := bare.SplitSpace(line)
	if len(parts) == 0 {
		return
	}
	cmd := parts[0]
	args := parts[1:]

	if fn, ok := commands[cmd]; ok {
		fn(args)
	} else {
		sys.Con.PutString("Unknown command: "+cmd+"\n", 0xFF0000FF)
	}
}

// --- Built-in command implementations ---

func CmdHelp(args []string) {
	sys.Con.PutString("Available commands:\n", 0x00FFFF00)
	for name := range commands {
		sys.Con.PutString(" - "+name+"\n", 0x00FFFF00)
	}
}

func CmdPrint(args []string) {
	if len(args) == 0 {
		sys.Con.PutString("Usage: print <text>\n", 0xFFAA00FF)
		return
	}
	text := bare.Join(args, ' ')
	sys.Con.PutString(text+"\n", 0xFFFFFFFF)
}

func CmdAdd(args []string) {
	if len(args) < 2 {
		sys.Con.PutString("Usage: add <num1> <num2>\n", 0xFFAA00FF)
		return
	}

	var a, b int32
	for i, s := range args[:2] {
		val, ok := bare.Atoi(s)
		if !ok {
			sys.Con.PutString("Invalid number: "+s+"\n", 0xFF0000FF)
			return
		}
		if i == 0 {
			a = val
		} else {
			b = val
		}
	}

	bare.PrintInt(a + b)
}
