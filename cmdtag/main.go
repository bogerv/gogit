package main

import (
	"fmt"
	"gitshell/cmdtag/config"
	"strings"

	"gitshell/cmds"
	"gitshell/colorlog"
)

func main() {
	var msgCh = make(chan string, 1)
	var endCh = make(chan struct{})
	conf := config.Init()

	go func(msgCh chan<- string) {
		for _, path := range conf.Paths {
			shell := new(cmds.Cmd)
			shell.Path = path
			// change path
			shell.Chdir()
			colorlog.Cyan(fmt.Sprintf(`change path to: "%s"`, path))

			// fetch latest commit for all branch
			shell.GitFetch()
			for _, branch := range conf.Branches {
				shell.CurrentBranch = branch

				tagName := fmt.Sprintf("tag-%s-%s", branch, conf.Version)
				basedBranch := fmt.Sprintf("origin/%s", branch)

				// remove local same tag
				shell.DeleteTag(tagName)

				// add tag
				shell.AddTag(tagName, basedBranch, conf.Message)

				// push tag
				shell.GitPushOrigin(tagName)

				if strings.EqualFold(branch, "pre") || strings.EqualFold(branch, "online") {
					// get latest commit id
					shell.GetCommitId()
					msgCh <- fmt.Sprintf("%s--%s--commit id-> %s", shell.Path, shell.CurrentBranch, shell.CommitId)
				} else {
					fmt.Println()
				}
			}
		}
		close(endCh)
	}(msgCh)

	for {
		select {
		case _, ok := <-endCh:
			if !ok {
				colorlog.Success("Finished...")
				return
			}
		case msg, ok := <-msgCh:
			if ok {
				colorlog.Success(msg)
			}
		}
	}
}
