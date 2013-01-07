package util

// Bytes to Word
func B2W(most byte, least byte) uint16 {
    return uint16(most)<<8 | uint16(least)
}

// Word to Bytes
func W2B(word uint16) (most byte, least byte) {
    most = byte((word >> 8) & 0xFF)
    least = byte(word & 0xFF)
    return
}
