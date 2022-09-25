package dudu

import (
	"os/exec"
	"strconv"
	"strings"
)

var duData = map[string]int{}

func ParseDuData(rawDuData []byte) map[string]int {
	duDataMap := map[string]int{}
	duSlice := strings.Fields(string(rawDuData))

	for s := 1; s < len(duSlice); s += 2 {
		size, _ := strconv.Atoi(duSlice[s-1])
		duDataMap[duSlice[s]] = size
	}

	return duDataMap
}

func Du(arg string) ([]byte, error) {
	// Run the du command
	if arg[len(arg)-1:] != "/" {
		arg = arg + "/"
	}
	duCmd := "du --summarize --exclude=proc --exclude=sys --exclude=dev --exclude=run --exclude=home " + arg + "*"
	du := exec.Command("/bin/sh", "-c", duCmd)
	rawDu, err := du.CombinedOutput()

	return rawDu, err
}
