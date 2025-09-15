// sys/console.go
package sys

import "unsafe"

type Console struct {
	Width      int
	Height     int
	CursorX    int
	CursorY    int
	FbWidth    int      // framebuffer width in pixels
	FbHeight   int      // framebuffer height in pixels
	FbBuffer   []uint32 // framebuffer slice
	CharWidth  int
	CharHeight int
}

// Global console instance
var Con Console

// Init initializes the console framebuffer
func (c *Console) Init() {
	c.CharWidth = 8
	c.CharHeight = 8

	c.FbWidth = 640
	c.FbHeight = 480

	info := fbInfo{
		width:      uint32(c.FbWidth),
		height:     uint32(c.FbHeight),
		virtWidth:  uint32(c.FbWidth),
		virtHeight: uint32(c.FbHeight),
		depth:      32,
	}
	addr := uintptr(unsafe.Pointer(&info)) &^ 0xF

	// Proper mailbox write/read
	mailboxWrite(uint32(addr), MAILBOX_CHANNEL_FB)
	resp := mailboxRead(MAILBOX_CHANNEL_FB)
	if resp == 0 || info.pointer == 0 || info.pitch == 0 {
		return // GPU did not respond
	}

	c.Width = c.FbWidth / c.CharWidth
	c.Height = c.FbHeight / c.CharHeight
	c.FbBuffer = unsafe.Slice((*uint32)(unsafe.Pointer(uintptr(info.pointer))), c.FbWidth*c.FbHeight)
	c.CursorX = 0
	c.CursorY = 0

	c.ClearScreen()
}

// ClearScreen fills framebuffer with black
func (c *Console) ClearScreen() {
	for i := range c.FbBuffer {
		c.FbBuffer[i] = 0x00000000
	}
	c.CursorX = 0
	c.CursorY = 0
}

// PutChar draws a single character at current cursor
func (c *Console) PutChar(ch byte, color uint32) {
	if ch == '\n' {
		c.CursorX = 0
		c.CursorY++
		if c.CursorY >= c.Height {
			c.scrollUp()
			c.CursorY--
		}
		return
	}

	FbPutChar(c.CursorX*c.CharWidth, c.CursorY*c.CharHeight, ch, color)
	c.CursorX++
	if c.CursorX >= c.Width {
		c.CursorX = 0
		c.CursorY++
		if c.CursorY >= c.Height {
			c.scrollUp()
			c.CursorY--
		}
	}
}

// PutString prints a string
func (c *Console) PutString(s string, color uint32) {
	for i := 0; i < len(s); i++ {
		c.PutChar(s[i], color)
	}
}

// scrollUp scrolls screen by one character row
func (c *Console) scrollUp() {
	// Move framebuffer up by CharHeight pixels
	for y := 0; y < c.FbHeight-c.CharHeight; y++ {
		for x := 0; x < c.FbWidth; x++ {
			c.FbBuffer[y*c.FbWidth+x] = c.FbBuffer[(y+c.CharHeight)*c.FbWidth+x]
		}
	}

	// Clear last row
	for y := c.FbHeight - c.CharHeight; y < c.FbHeight; y++ {
		for x := 0; x < c.FbWidth; x++ {
			c.FbBuffer[y*c.FbWidth+x] = 0x00000000
		}
	}
}
