package polylib

import (
	"fmt"
	"io"
)

func Polyfmt(s, colorHex string) {
	fmt.Printf("%%{F%s}%s%%{F-}", colorHex, s)
}

// Does not flush the writer.
func PolyfmtBufWr(s, colorHex string, wr io.Writer) {
	v := fmt.Sprintf("%%{F%s}%s%%{F-}", colorHex, s)
	wr.Write([]byte(v))
}

// Red hex. Not 100% sure when it needs newline terminator
// if something crashes
func Polyerr(err error) {
	Polyfmt(err.Error(), "#F54242")
}

// To test hex colors in equivalent rgb. Prints to stdout.
func DebugColors(s string, r, g, b int) {
	fmt.Printf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, s)
}

func getHexes(n int) (int, int) {
	return n / 16, n % 16
}

func getByte(n int) byte {
	if n < 10 {
		return 48 + byte(n)
	}
	switch n {
	case 10:
		return 'A'
	case 11:
		return 'B'
	case 12:
		return 'C'
	case 13:
		return 'D'
	case 14:
		return 'E'
	case 15:
		return 'F'
	}
	return 'A'
}

// Includes the '#' at the beginning.
// Assumes 255 max.
func RgbToHex(r, g, b int) string {
	colors := [3]int{r, g, b}
	bytes := []byte("#000000")
	idx := 1
	for _, c := range colors {
		div, rem := getHexes(c)
		tensLetter := getByte(div)
		letter := getByte(rem)
		bytes[idx] = tensLetter
		bytes[idx+1] = letter
		idx += 2
	}
	return string(bytes)
}
