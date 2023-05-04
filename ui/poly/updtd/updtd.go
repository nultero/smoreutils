package main

import (
	"errors"
	"fmt"
	"os"
	"polylib"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	chErr      = errors.New("inotify channel failed somehow")
	hasUpdated = false
	bellIcon   = ""
	writeOp    = fsnotify.Op(2)
)

func printBell() {
	// TODOO change this to adapt colors dynamically
	polylib.Polyfmt(bellIcon, "#f4ff04")
	fmt.Print("\n")
}

func watchUpdate(w *fsnotify.Watcher, fpath string) {
	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				polylib.Polyerr(chErr)
				continue
			}

			if ev.Has(writeOp) {
				if isUpdated(fpath) {
					fmt.Print(" \n")
				} else {
					printBell()
				}
			}

		case err, ok := <-w.Errors:
			if !ok {
				polylib.Polyerr(chErr)
				continue
			}
			polylib.Polyerr(err)
		}
	}
}

func isUpdated(fpath string) bool {
	now := time.Now()
	mon := now.Month().String()[:3]
	day := now.Day()
	currDate := fmt.Sprintf("%d %s %d", now.Year(), mon, day)

	bytes, err := os.ReadFile(fpath)
	if err != nil {
		polylib.Polyerr(err)
	}
	lines := strings.Split(string(bytes), "\n")
	return currDate == lines[0]
}

func eternalSleep(fpath string) {
	// TODO maybe sleep this until
	// the "next" day, in case you
	// burn midnight oil
	for {
		if !isUpdated(fpath) {
			printBell()
		}
		time.Sleep(1 * time.Hour)
	}
}

// not Go's init:
// polybar (as of my last version) will not show anything
// if on startup, there is no output -- such as in the case
// of reboooting after already updating that day
func doInit() {
	fmt.Print("x\n")
	time.Sleep(400 * time.Millisecond)
	fmt.Print(" \n")
}

func main() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		polylib.Polyerr(err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		polylib.Polyerr(err)
	}
	updtPath := fmt.Sprintf("%s/.updt_last_run.txt", home)

	go watchUpdate(w, updtPath)

	err = w.Add(updtPath)
	if err != nil {
		polylib.Polyerr(err)
	}

	doInit()
	eternalSleep(updtPath)
}
