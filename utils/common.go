package utils

import (
	"strings"
)

var umlautsMapping = map[string]string{
	"ä": "ae",
	"ü": "ue",
	"ö": "oe",
	"ß": "sz",
}

// SanitizeUmlauts replaces German umlauts with their two-letter equivalents.
func SanitizeUmlauts(message string) string {
	for key, value := range umlautsMapping {
		message = strings.ReplaceAll(message, key, value)
	}
	return message
}

// Byte2Hex converts a byte to a hexadecimal string representation.
// This function's logic seems unusual for a typical byte to hex conversion.
// It appears to be converting decimal digits of the byte into a single character
// based on a custom encoding (shifting high nibble by 4 and adding low nibble).
// This might be specific to the LED board communication protocol.
func Byte2Hex(b byte) string {
	high := b / 10
	low := b % 10
	return string(rune((high << 4) + low))
}
