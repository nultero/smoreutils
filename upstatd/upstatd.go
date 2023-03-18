package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// TODO old regex from conf -> #[a-fA-F0-9]{6}$
// TODOOOO make smoreutils opts to hide for desktop w/ no battery

var (
	engNow         = "/sys/class/power_supply/BAT0/energy_now"
	engFul         = "/sys/class/power_supply/BAT0/energy_full"
	lightningBolt  = "⚡"
	lightningGreen = "#11F71D" // TODO darken hex based on bat level
	textMagenta    = "#AD003D"
	thresholds     = []float64{92.0, 83.0, 70.0, 50.0, 24.0, -0.1}
	batChars       = []string{"", "", "", "", "", "!"}
	max            = 97.0
)

func getPercent(buf []byte, fd *os.File) (float64, error) {
	_, err := fd.Seek(0, 0)
	if err != nil {
		polyerr(err)
		return 0.0, err
	}

	n, err := fd.Read(buf)
	if err != nil {
		polyerr(err)
		return 0.0, err
	}
	s := string(buf[:n-1]) //slice off newline
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		polyerr(err)
	}
	return f, err
}

// TODO read from polybar confs
func main() {

	enowf, err := os.Open(engNow)
	if err != nil {
		polyerr(err)
		fmt.Print("\n")
		os.Exit(1)
	}
	efulf, err := os.Open(engNow)
	if err != nil {
		polyerr(err)
		fmt.Print("\n")
		os.Exit(1)
	}

	buf := make([]byte, 20)
	enow, eful, batPc := 0.0, 0.0, 0.0

	for {
		enow, err = getPercent(buf, enowf)
		if err != nil {
			goto endOfLoop
		}
		eful, err = getPercent(buf, efulf)
		if err != nil {
			goto endOfLoop
		}

		polyfmt("BAT ", textMagenta)

		batPc = (enow / eful) * 100.0

		if batPc > max {
			polyfmt(lightningBolt, lightningGreen)
			goto endOfLoop
		}

		for i := 0; i < 6; i++ {
			if batPc > thresholds[i] {
				polyfmt(batChars[i], lightningGreen)
				break
			}
		}

	endOfLoop:
		fmt.Print("\n")
		time.Sleep(1 * time.Second)
	}
}
