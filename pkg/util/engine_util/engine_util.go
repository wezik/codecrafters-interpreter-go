package util

import "strings"

func IsWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '	' || b == '\n'
}

func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func IsFirstDot(b byte, buffer string) bool {
	return b == '.' && !strings.Contains(buffer, ".")
}

func IsAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
