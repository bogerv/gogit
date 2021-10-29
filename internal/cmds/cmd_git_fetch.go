package cmds

import (
	"fmt"
	"os/exec"
)

// GitFetch git fetch
func (slf *Cmd) GitFetch() {
	cmd := exec.Command("git", "fetch")
	fmt.Println("exec command:: " + fmt.Sprintf("git fetch"))

	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("git fetch error: %v", err))
	}
}
