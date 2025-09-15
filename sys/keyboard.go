// sys/keyboard.go
package sys

// --- USB keyboard registers (simplified for BCM2835/2837) ---
// For real implementation, you'd use USB host controller (EHCI) registers.
// We'll mock memory-mapped buffer for demo purposes.
var usbKeyboardBuffer = [8]byte{} // HID report: 8 bytes
var lastKey byte = 0

// ReadKey polls the USB keyboard and returns a single ASCII character
// Returns 0 if no key is pressed
func ReadKey() byte {
	// Poll USB keyboard HID report (mocked here)
	if usbKeyboardBuffer[2] != 0 { // keycode in byte 2
		key := hidToAscii(usbKeyboardBuffer[2])
		usbKeyboardBuffer[2] = 0 // consume
		lastKey = key
		return key
	}
	return 0
}

// Convert HID keycode (simplified) to ASCII
func hidToAscii(hid byte) byte {
	// Only basic keys: letters a-z, 1-9, 0, Enter, Backspace
	switch hid {
	case 0x04:
		return 'a'
	case 0x05:
		return 'b'
	case 0x06:
		return 'c'
	case 0x07:
		return 'd'
	case 0x08:
		return 'e'
	case 0x09:
		return 'f'
	case 0x0A:
		return 'g'
	case 0x0B:
		return 'h'
	case 0x0C:
		return 'i'
	case 0x0D:
		return 'j'
	case 0x0E:
		return 'k'
	case 0x0F:
		return 'l'
	case 0x10:
		return 'm'
	case 0x11:
		return 'n'
	case 0x12:
		return 'o'
	case 0x13:
		return 'p'
	case 0x14:
		return 'q'
	case 0x15:
		return 'r'
	case 0x16:
		return 's'
	case 0x17:
		return 't'
	case 0x18:
		return 'u'
	case 0x19:
		return 'v'
	case 0x1A:
		return 'w'
	case 0x1B:
		return 'x'
	case 0x1C:
		return 'y'
	case 0x1D:
		return 'z'
	case 0x1E:
		return '1'
	case 0x1F:
		return '2'
	case 0x20:
		return '3'
	case 0x21:
		return '4'
	case 0x22:
		return '5'
	case 0x23:
		return '6'
	case 0x24:
		return '7'
	case 0x25:
		return '8'
	case 0x26:
		return '9'
	case 0x27:
		return '0'
	case 0x28:
		return '\n' // Enter
	case 0x2A:
		return '\b' // Backspace
	default:
		return 0
	}
}

// --- Testing helper ---
func SimulateKeyPress(hid byte) {
	usbKeyboardBuffer[2] = hid
}
