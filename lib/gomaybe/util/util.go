package util

// Bytes to Word
func B2W(most byte, least byte) uint16 {
    return uint16(most)<<8 | uint16(least)
}
