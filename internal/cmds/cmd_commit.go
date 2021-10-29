package cmds

import (
	"bytes"
	"fmt"
	"gitshell/pkg/colorlog"
	"os/exec"
	"strings"
)

// Add exec git add . command
func (slf *Cmd) Commit(msg string) {
	args := []string{"commit", "-m", msg}
	cmd := exec.Command("git", args...)
	fmt.Println("exec command:: git " + strings.Join(args, " "))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		colorlog.Error(fmt.Sprintf("git commit -m '' error: %v, detail: %s", err, stderr.String()))
	}
}
