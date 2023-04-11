package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func redErr(err error) {
	fmt.Printf("\x1b[31merr:\x1b[0m %s\n", err.Error())
	os.Exit(1)
}

func getDateStr() string {
	t := time.Now()
	mon := t.Month().String()[:3]
	s := fmt.Sprintf(
		"%d %s %d\n",
		t.Year(), mon, t.Day(),
	)
	return s
}

func getDistro() string {
	bytes, err := os.ReadFile("/etc/os-release")
	if err != nil {
		redErr(err)
	}

	err = errors.New("/etc/os-release does not have line 'ID=...' or this line is broken")
	var release *string

	lines := strings.Split(string(bytes), "\n")
	for _, ln := range lines {
		if ln[:3] != "ID=" {
			continue
		}

		spl := strings.Split(ln, "=")
		if len(spl) < 2 {
			redErr(err)
		}
		release = &spl[1]
		break
	}

	if release == nil {
		redErr(err)
	}

	return *release
}

func setLastUpdt(datestr string) {
	home, err := os.UserHomeDir()
	if err != nil {
		redErr(err)
	}

	fpath := home + "/.updt_last_run.txt"
	bytes, err := os.ReadFile(fpath)
	if err != nil {
		redErr(err)
	}
	str := string(bytes)
	lines := strings.Split(str, "\n")
	lines[0] = datestr
	str = strings.Join(lines, "\n")
	bytes = []byte(str)
	err = os.WriteFile(fpath, bytes, 0744)
	if err != nil {
		redErr(err)
	}
	fmt.Println("\x1b[38;2;0;200;0mUPDATE SUCCESSFUL\x1b[0m")
}

func getCmds(distro string) []*exec.Cmd {
	cmds := []string{}
	switch distro {
	case "arch":
		cmds = append(cmds, "sudo pacman -Syu")

	case "debian", "ubuntu":
		cmds = append(
			cmds,
			"sudo apt update",
			"sudo apt upgrade",
		)

	case "debug":
	case "fedora":
	}
	execs := []*exec.Cmd{}
	for _, c := range cmds {
		strs := strings.Split(c, " ")
		cmd := exec.Command(strs[0], strs[1:]...)
		execs = append(execs, cmd)
	}
	return execs
}

func displayCmd(c *exec.Cmd) {
	s := strings.Join(c.Args, " ")
	fmt.Printf("running: \x1b[33m%s\x1b[0m\n", s)
}

func main() {

	// TODO -> already updated
	date := getDateStr()
	distro := getDistro()

	for _, arg := range os.Args[1:] {
		if arg == "-g" {
			distro = "debug"
		}
	}

	cmds := getCmds(distro)
	for _, cmd := range cmds {
		displayCmd(cmd)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			redErr(err)
		}
	}

	setLastUpdt(date)
}
