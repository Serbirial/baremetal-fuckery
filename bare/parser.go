// bare/parser.go
package bare

import (
	"fmt"
	"strings"
)

// Simple assembler-like parser for our VM
// Format: "PUSH 10", "ADD", "PRINT"
func ParseSource(src string) ([]byte, error) {
	lines := strings.Split(src, "\n")
	bytecode := make([]byte, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		parts := strings.Fields(line)
		switch parts[0] {
		case "PUSH":
			bytecode = append(bytecode, 0x01)
			var val int32
			fmt.Sscanf(parts[1], "%d", &val)
			bytecode = append(bytecode, byte(val), byte(val>>8), byte(val>>16), byte(val>>24))
		case "ADD":
			bytecode = append(bytecode, 0x02)
		case "SUB":
			bytecode = append(bytecode, 0x03)
		case "PRINT":
			bytecode = append(bytecode, 0x04)
		default:
			return nil, fmt.Errorf("unknown instruction: %s", parts[0])
		}
	}

	return bytecode, nil
}
