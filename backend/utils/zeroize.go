package utils

// Zeroize overwrites a byte slice with zeros to prevent key material lingering in memory.
func Zeroize(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
