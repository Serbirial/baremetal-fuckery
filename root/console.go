// root/console.go
package root

import "sys"

// REPLConsole wraps the sys console for REPL usage
type REPLConsole struct {
	prompt    string
	buffer    []byte
	cursor    int
	history   [][]byte
	histIndex int
}

// Global REPL console instance
var Con REPLConsole

// Init initializes the REPL console with a prompt
func (c *REPLConsole) Init(prompt string) {
	c.prompt = prompt
	c.buffer = make([]byte, 0, 256)
	c.cursor = 0
	c.history = make([][]byte, 0)
	c.histIndex = -1

	// Draw initial prompt
	sys.Con.PutString(c.prompt, 0x00FF00FF)
}

// PutChar adds a character to the REPL buffer and draws it
func (c *REPLConsole) PutChar(ch byte) {
	// Print to framebuffer
	sys.Con.PutChar(ch, 0xFFFFFFFF)

	// Append to buffer
	c.buffer = append(c.buffer, ch)
	c.cursor++
}

// Backspace removes the last character from buffer and framebuffer
func (c *REPLConsole) Backspace() {
	if c.cursor == 0 {
		return
	}
	c.cursor--
	c.buffer = c.buffer[:len(c.buffer)-1]

	// Move cursor back visually and overwrite with space
	if sys.Con.CursorX > 0 {
		sys.Con.CursorX--
		sys.Con.PutChar(' ', 0x00000000)
		sys.Con.CursorX--
	}
}

// ExecuteBuffer sends the buffer to command parser and clears it
func (c *REPLConsole) ExecuteBuffer() {
	line := string(c.buffer)
	c.history = append(c.history, append([]byte(nil), c.buffer...)) // copy buffer
	c.histIndex = len(c.history)
	c.buffer = c.buffer[:0]
	c.cursor = 0

	// Newline in console
	sys.Con.PutChar('\n', 0xFFFFFFFF)

	// Execute command
	ExecuteCommand(line)

	// Print prompt again
	sys.Con.PutString(c.prompt, 0x00FF00FF)
}

// HandleKey processes basic keys (ASCII only for now)
func (c *REPLConsole) HandleKey(ch byte) {
	switch ch {
	case '\r', '\n':
		c.ExecuteBuffer()
	case 8, 127: // Backspace
		c.Backspace()
	default:
		c.PutChar(ch)
	}
}
