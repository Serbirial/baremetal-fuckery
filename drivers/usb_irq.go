// drivers/usb_irq.go
package drivers

import "sys"

// Dummy IRQ number for USB (replace with actual vector in vector table)
const IRQ_USB = 9

// USB interrupt handler
func UsbIRQHandler() {
	// Check if port has pending interrupt
	status := USB.CmdSts.Read()
	if status == 0 {
		return
	}

	// For simplicity, we assume only keyboard endpoint IN has data
	// Copy incoming keycode to Keyboard.report[2]
	key := uint8(status & 0xFF) // low byte: HID keycode
	if key != 0 {
		Keyboard.report[2] = key
	}

	// Clear interrupt
	USB.CmdSts.Write(status)
}

// Install handler
func InstallUSBIRQ() {
	sys.Con.PutString("Installing USB IRQ handler...\n", 0x00FF00FF)
	sys.RegisterIRQ(IRQ_USB, UsbIRQHandler)
	sys.EnableIRQ(IRQ_USB)
}
