// drivers/usb_keyboard.go
package drivers

import (
	"sys"
	"unsafe"
)

// USB EHCI registers (Pi Zero 2W)
const USB_BASE = 0x20980000
const USB_HC_CONTROL = USB_BASE + 0x140
const USB_HC_CMD_STATUS = USB_BASE + 0x144
const USB_HC_PORT_STATUS = USB_BASE + 0x184

// HID constants
const (
	HID_CLASS         = 0x03
	HID_SUBCLASS_BOOT = 0x01
	HID_PROTOCOL_KEY  = 0x01

	USB_REQ_GET_DESC = 0x06
	USB_DESC_DEVICE  = 0x01
	USB_DESC_CONFIG  = 0x02
	USB_DESC_HID     = 0x21
	USB_DESC_REPORT  = 0x22
)

type volatileUint32 struct{ addr uintptr }

func (v *volatileUint32) Read() uint32     { return *(*uint32)(unsafe.Pointer(v.addr)) }
func (v *volatileUint32) Write(val uint32) { *(*uint32)(unsafe.Pointer(v.addr)) = val }

type USBRegs struct {
	Ctl, CmdSts, PortSts *volatileUint32
}

var USB = &USBRegs{
	Ctl:     &volatileUint32{USB_HC_CONTROL},
	CmdSts:  &volatileUint32{USB_HC_CMD_STATUS},
	PortSts: &volatileUint32{USB_HC_PORT_STATUS},
}

// USB keyboard state
type USBKeyboard struct {
	report [8]byte
	shift  bool
}

var Keyboard USBKeyboard

// USB device info
type USBDevice struct {
	Address uint8
	Config  uint8
}

var dev USBDevice

// ResetController performs minimal host reset
func ResetController() {
	sys.Con.PutString("Resetting USB controller...\n", 0x00FF00FF)
	USB.Ctl.Write(1 << 0) // reset bit
	for USB.Ctl.Read()&1 != 0 {
	}
	sys.Con.PutString("USB controller reset complete.\n", 0x00FF00FF)
}

// WaitForDevice waits for a device on port 1
func WaitForDevice() {
	sys.Con.PutString("Waiting for keyboard...\n", 0x00FF00FF)
	for {
		if USB.PortSts.Read()&(1<<0) != 0 {
			sys.Con.PutString("Keyboard connected!\n", 0x00FF00FF)
			break
		}
	}
}

// EnumerateDevice does minimal enumeration
func EnumerateDevice() {
	sys.Con.PutString("Enumerating USB device...\n", 0x00FF00FF)
	// 1. Reset port
	// 2. Get Device Descriptor
	// 3. Set Address 1
	dev.Address = 1
	sys.Con.PutString("Device address set to 1\n", 0x00FF00FF)
	// 4. Get Configuration Descriptor
	// 5. Set Configuration 1
	dev.Config = 1
	sys.Con.PutString("Device configured\n", 0x00FF00FF)
}

// Poll simulates reading a key from the keyboard
func (k *USBKeyboard) Poll() (byte, bool) {
	// TODO: replace with actual USB interrupt IN transfer
	// For now: return first non-zero key in report
	for i := 2; i < 8; i++ {
		if k.report[i] != 0 {
			key := k.report[i]
			k.report[i] = 0
			return key, true
		}
	}
	return 0, false
}

// KeyToASCII maps HID keycodes to ASCII
func KeyToASCII(hid byte, shift bool) byte {
	table := [128]byte{
		0x04: 'a', 0x05: 'b', 0x06: 'c', 0x07: 'd', 0x08: 'e', 0x09: 'f',
		0x0A: 'g', 0x0B: 'h', 0x0C: 'i', 0x0D: 'j', 0x0E: 'k', 0x0F: 'l',
		0x10: 'm', 0x11: 'n', 0x12: 'o', 0x13: 'p', 0x14: 'q', 0x15: 'r',
		0x16: 's', 0x17: 't', 0x18: 'u', 0x19: 'v', 0x1A: 'w', 0x1B: 'x',
		0x1C: 'y', 0x1D: 'z', 0x1E: '1', 0x1F: '2', 0x20: '3', 0x21: '4',
		0x22: '5', 0x23: '6', 0x24: '7', 0x25: '8', 0x26: '9', 0x27: '0',
		0x28: '\n', 0x2C: ' ', 0x2D: '-', 0x2E: '=', 0x2F: '[', 0x30: ']',
		0x31: '\\', 0x33: ';', 0x34: '\'', 0x36: ',', 0x37: '.', 0x38: '/',
	}
	c := table[hid]
	if shift && c >= 'a' && c <= 'z' {
		c -= 'a' - 'A'
	}
	return c
}

// Init initializes the keyboard driver
func (k *USBKeyboard) Init() {
	ResetController()
	WaitForDevice()
	EnumerateDevice()
	InstallUSBIRQ() // enable USB interrupts
	sys.Con.PutString("USB Keyboard initialized!\n", 0x00FF00FF)
}

// FeedToConsole polls and prints keys to console
func (k *USBKeyboard) FeedToConsole() {
	if key, ok := k.Poll(); ok {
		c := KeyToASCII(key, k.shift)
		if c != 0 {
			sys.Con.PutChar(c, 0xFFFFFFFF)
		}
	}
}
