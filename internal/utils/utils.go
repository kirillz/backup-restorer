package utils

import (
	"os"
	"os/exec"
	"runtime"
)

func ClearTerminal() {
	
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
