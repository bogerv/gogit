package cmds

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitFetch git fetch
func (slf *Cmd) GitCheckout(branch string) {
	args := []string{"checkout", branch}
	cmd := exec.Command("git", args...)
	fmt.Println("exec command:: git", strings.Join(args, " "))

	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("git checkout error: %v", err))
	}
}

