package cmds

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitPushOrigin 向远端推送最新的分支或者Tag标签
func (slf *Cmd) GitPull() {
	args := []string{"pull"}
	cmd := exec.Command("git", args...)
	fmt.Println("exec command:: git "+ strings.Join(args, " "))

	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("git pull error: %v", err))
	}
}
