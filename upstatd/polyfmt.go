package main

import "fmt"

func polyfmt(s, colorHex string) {
	fmt.Printf("%%{F%s}%s%%{F-}", colorHex, s)
}

// Red hex. Not 100% sure if it needs newline terminator
// if something crashes
func polyerr(err error) {
	polyfmt(err.Error(), "#F54242")
}
