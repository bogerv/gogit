package cmds

import (
	"bytes"
	"fmt"
	"gitshell/pkg/colorlog"
	"os/exec"
	"strings"
)

// Add exec git add . command
func (slf *Cmd) Add(files ...string) {
	if len(files) <= 0 {
		args := []string{"add", "."}
		cmd := exec.Command("git", args...)
		fmt.Println("exec command:: git " + strings.Join(args, " "))

		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			colorlog.Error(fmt.Sprintf("git add . error: %v, detail: %s", err, stderr.String()))
			panic(err)
		}
	}
}
