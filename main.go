package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"git.lupinelab.co.uk/lupinelab/dudu/cmd/dudu"
)

func checkSudo() string {
	stdout, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(stdout)
}

func main() {
	if checkSudo() != "root\n" {
		fmt.Println("Please run as elevated user")
		return
	}
	dudu.Execute()
}
