package utils

func Zeroize(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
