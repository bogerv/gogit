package cmds

import (
	"bytes"
	"fmt"
	"gitshell/pkg/colorlog"
	"os/exec"
	"strings"
)

// GitPushOrigin 向远端推送最新的分支或者Tag标签
func (slf *Cmd) GitPushOrigin(tagOrBranch string) {
	args := []string{"push", "origin", tagOrBranch}
	cmd := exec.Command("git", args...)
	fmt.Println("exec command:: git " + strings.Join(args, " "))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		colorlog.Error(fmt.Sprintf("git push origin error: %v, detail: %s", err, stderr.String()))
		panic(err)
	}
}
