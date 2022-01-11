package cmds

import (
	"fmt"
	"gitshell/pkg/colorlog"
	"os/exec"
)

// GetLog exec git command `git log -n 1 --pretty=format:"%H"`
func (slf *Cmd) GetLog(branch string) {
	//cmsArgs := []string{"log", "-n", "1", `--pretty=format:"%H"`}
	cmsArgs := []string{"log", "-n", "1", `--pretty=format:%H`}
	cmd := exec.Command("git", cmsArgs...)
	fmt.Println("exec command:: git", "log", "-n", "1", `--pretty=format:%H`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("git log error: %v", err))
	}

	colorlog.Green(fmt.Sprintf("%s commit id: %s\n", branch, output))
}

// GetCommitId exec git command `git rev-parse HEAD`
func (slf *Cmd) GetCommitId() {
	//cmdArgs := []string{"rev-parse", "HEAD"}
	cmdArgs := []string{"rev-parse", "origin/" + slf.CurrentBranch}
	cmd := exec.Command("git", cmdArgs...)
	//fmt.Println("exec command:: git", strings.Join(cmdArgs, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("git rev-parse HEAD error: %v", err))
	}

	slf.CommitId = string(output)
}
