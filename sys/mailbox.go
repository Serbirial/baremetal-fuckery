// sys/mailbox.go
package sys

import "unsafe"

// mailboxWrite writes a message to the specified channel
func mailboxWrite(msg uint32, channel uint32) {
	// The lower 4 bits are reserved for the channel
	if channel > 0xF {
		return
	}

	// Construct the value: lower 4 bits = channel
	val := (msg & 0xFFFFFFF0) | (channel & 0xF)

	// Pointer to mailbox
	mailbox := uintptr(MAILBOX_BASE)

	// Wait until mailbox is not full
	for *(*uint32)(unsafe.Pointer(mailbox + MAILBOX_STATUS))&MAILBOX_FULL != 0 {
	}

	// Write the value
	*(*uint32)(unsafe.Pointer(mailbox + MAILBOX_WRITE)) = val
}

// mailboxRead reads a message from the specified channel
func mailboxRead(channel uint32) uint32 {
	if channel > 0xF {
		return 0
	}

	mailbox := uintptr(MAILBOX_BASE)

	for {
		// Wait until mailbox is not empty
		for *(*uint32)(unsafe.Pointer(mailbox + MAILBOX_STATUS))&MAILBOX_EMPTY != 0 {
		}

		// Read the value
		val := *(*uint32)(unsafe.Pointer(mailbox + MAILBOX_READ))

		// Lower 4 bits = channel
		if val&0xF == channel {
			return val & 0xFFFFFFF0
		}
	}
}
