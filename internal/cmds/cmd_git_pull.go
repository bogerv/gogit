package cmds

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitPull 拉取分支最新内容
func (slf *Cmd) GitPull() {
	args := []string{"pull"}
	cmd := exec.Command("git", args...)
	fmt.Println("exec command:: git " + strings.Join(args, " "))

	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("git pull error: %v", err))
	}
}
