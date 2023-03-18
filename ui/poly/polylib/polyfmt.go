package polylib

import "fmt"

func Polyfmt(s, colorHex string) {
	fmt.Printf("%%{F%s}%s%%{F-}", colorHex, s)
}

// Red hex. Not 100% sure when it needs newline terminator
// if something crashes
func Polyerr(err error) {
	Polyfmt(err.Error(), "#F54242")
}
