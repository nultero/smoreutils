package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"polylib"
	"runtime"
	"strings"
	"time"
)

// 255, 102, 0 -- orange in rgb
// ff6600

// #ffff00
//rgb(255, 255, 0) -- yellow rgb

const icon = ""

var newline = []byte{'\n'}

type cpu struct {
	id,
	work, idle int
}

func (c cpu) usage(prev cpu) float32 {
	wrk := float32(c.work - prev.work)
	idle := float32(c.idle - prev.idle)
	if idle == 0.0 {
		return 0.0
	}
	return (wrk / (wrk + idle))
}

func main() {
	procStat, err := os.Open("/proc/stat")
	if err != nil {
		polylib.Polyerr(err)
		os.Exit(1)
	}

	numCpus := runtime.NumCPU()
	buf := make([]byte, 1000)
	cpus := make([]cpu, numCpus)
	sats := make([]float32, numCpus)
	stdout := bufio.NewWriter(os.Stdout)

	for {
		procStat.Seek(0, 0)

		_, err := procStat.Read(buf)
		if err != nil {
			polylib.Polyerr(err)
			os.Exit(1)
		}
		newCpus := parseBuf(buf, numCpus)

	miniloop:
		for idx := 0; idx < numCpus; idx++ {
			prev := cpus[idx]
			if prev.id == 0 {
				cpus = newCpus
				break miniloop
			}
			usage := newCpus[idx].usage(prev)
			cpus[idx] = newCpus[idx]
			sats[idx] = usage
		}

		for _, sat := range sats {
			hex := fmtSaturation(sat)
			polylib.PolyfmtBufWr(icon, hex, stdout)
		}

		stdout.Write(newline)
		stdout.Flush()
		time.Sleep(250 * time.Millisecond)
	}
}

const ( //rgb vals
	orangeRed   float32 = 255.0
	orangeGreen float32 = 102.0
	orangeBlue  float32 = 0.0
	threshold   float32 = 0.25 // the shift from orange to red, at 75% use
)

func fmtSaturation(sat float32) string {
	r, g := 0, 0
	thresh := 1.0 - threshold
	if sat < thresh {
		scalePercent := 1.0 / thresh
		r = int(sat * scalePercent * 1.18 * orangeRed)
		if r > 255 {
			r = 255
		}
		g = int(sat * scalePercent * orangeGreen)

	} else {
		r = 255
		chunk := 1.0 - thresh
		diff := 1.0 - sat
		percent := diff / chunk
		g = int(percent)
	}

	return polylib.RgbToHex(r, g, 0)
}

type serials struct {
	user, nice, system, idle, iowait,
	irq, softirq, steal, guest, guestNice,
	id int
}

func parseBuf(buf []byte, numCpus int) []cpu {
	cpus := make([]cpu, numCpus)
	split := strings.Split(string(buf), "\n")

	t := serials{
		user:      0,
		nice:      0,
		system:    0,
		idle:      0,
		iowait:    0,
		irq:       0,
		softirq:   0,
		steal:     0,
		guest:     0,
		guestNice: 0,
		id:        0,
	}

	for idx := 1; idx <= numCpus; idx++ {
		line := split[idx][3:]
		n, err := fmt.Sscanf(
			line, "%d %d %d %d %d %d %d %d %d %d %d",
			&t.id, &t.user, &t.nice, &t.system, &t.idle, &t.iowait,
			&t.irq, &t.softirq, &t.steal, &t.guest, &t.guestNice,
		)
		c := cpu{
			id:   t.id + 1,
			work: 0,
			idle: 0,
		}

		if n != 11 {
			s := fmt.Sprintf("cpu%d weird values", idx-1)
			polylib.Polyerr(errors.New(s))
		} else if err != nil {
			polylib.Polyerr(err)
		}

		c.idle = t.idle + t.iowait
		c.work = t.user + t.nice + t.system + t.irq + t.softirq
		cpus[idx-1] = c
	}

	return cpus
}
