package dudu

import (
	// "fmt"
	"os/exec"
)

func Du(arg string) (string, error) {
	duCmd := "du -s --exclude=/{proc,sys,dev,run,home} " + arg + "*"
	du := exec.Command("/bin/sh", "-c", string(duCmd))
	thisrun, err := du.CombinedOutput()

	return string(thisrun), err
}
