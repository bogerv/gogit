package cmds

import (
	"fmt"
	"github.com/go-errors/errors"
	"os/exec"
)

// AddTag add local tag
func (slf *Cmd) AddTag(tagName, basedBranch, message string) {
	defer CatchPanic()

	cmd := exec.Command("git", "tag", "-a", tagName, basedBranch, "-m", message)
	fmt.Println("exec command:: " + fmt.Sprintf("git tag -a %s %s -m \"%s\"", tagName, basedBranch, message))

	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("git tag -a error: %v", err))
	}
}

// DeleteTag delete local same tag
func (slf *Cmd) DeleteTag(tagName string) {
	defer CatchPanic()

	cmd := exec.Command("git", "tag", "-d", tagName)
	fmt.Println("exec command:: " + fmt.Sprintf("git tag -d %s", tagName))

	err := cmd.Run()
	if err != nil {
		panic(errors.New("git tag -d error: " + fmt.Sprintf("%v", err)))
	}
}

// PushOriginDelete git push origin --delete tag
func (slf *Cmd) PushOriginDelete(tagName string) {
	defer CatchPanic()

	cmd := exec.Command("git", "push", "origin", "--delete", "tag", tagName)
	fmt.Println("exec command:: " + fmt.Sprintf("git push origin --delete tag %s", tagName))

	err := cmd.Run()
	if err != nil {
		panic(errors.New("git push origin --delete tag error: " + fmt.Sprintf("%v", err)))
	}
}

// PushAllTags git push origin --tags
func (slf *Cmd) PushAllTags() {
	defer CatchPanic()

	cmd := exec.Command("git", "push", "origin", "--tags")
	fmt.Println("exec command:: git push origin --tags")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}