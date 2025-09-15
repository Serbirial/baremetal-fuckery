// bare/stdlib.go
package bare

import "sys"

// ----------------------------
// Printing integers
// ----------------------------
func PrintInt(val int32) {
	buf := [12]byte{} // enough for -2147483648
	i := 0
	negative := false
	if val < 0 {
		negative = true
		val = -val
	}
	if val == 0 {
		buf[i] = '0'
		i++
	}
	for val > 0 {
		buf[i] = byte('0' + val%10)
		val /= 10
		i++
	}
	if negative {
		buf[i] = '-'
		i++
	}
	// Reverse the buffer
	for j := 0; j < i/2; j++ {
		buf[j], buf[i-1-j] = buf[i-1-j], buf[j]
	}
	// Print to console
	sys.Con.PutString(string(buf[:i]), 0xFFFFFFFF)
	sys.Con.PutChar('\n', 0xFFFFFFFF)
}

// ----------------------------
// Minimal string helpers
// ----------------------------

// Trim removes leading/trailing spaces or tabs
func Trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// SplitSpace splits a string by spaces/tabs
func SplitSpace(s string) []string {
	var parts []string
	i := 0
	for i < len(s) {
		// skip spaces
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}
		start := i
		for i < len(s) && s[i] != ' ' && s[i] != '\t' {
			i++
		}
		if start < i {
			parts = append(parts, s[start:i])
		}
	}
	return parts
}

// Join combines string slices with a separator
func Join(parts []string, sep byte) string {
	var b []byte
	for i, s := range parts {
		for j := 0; j < len(s); j++ {
			b = append(b, s[j])
		}
		if i < len(parts)-1 {
			b = append(b, sep)
		}
	}
	return string(b)
}

// Atoi converts string to int32, returns false if invalid
func Atoi(s string) (int32, bool) {
	var val int32
	neg := false
	i := 0
	if len(s) > 0 && s[0] == '-' {
		neg = true
		i = 1
	}
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		val = val*10 + int32(c-'0')
	}
	if neg {
		val = -val
	}
	return val, true
}
