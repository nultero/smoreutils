package polylib

import "testing"

func Test_hex(t *testing.T) {
	s := RgbToHex(250, 91, 61)
	w := "#FA5B3D"
	if s != w {
		t.Errorf("got %s, want %s", s, w)
	}
}
