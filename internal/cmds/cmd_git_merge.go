package cmds

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GitMerge exec command: git merge [branch]
func (slf *Cmd) GitMerge(branch string) {
	args := []string{"merge", branch}
	cmd := exec.Command("git", args...)
	fmt.Println("exec command:: git " + strings.Join(args, " "))

	//var out bytes.Buffer
	//cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		panic("git merge error: " + fmt.Sprint(err) + ": " + stderr.String())
	}
}
