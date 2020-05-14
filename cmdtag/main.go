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
	var endCh = make(chan bool)
	conf := config.Init()

	go func(ch chan<- string) {
		for _, path := range conf.Paths {
			shell := new(cmds.Cmd)
			shell.Path = path
			// change path
			shell.Chdir()
			colorlog.Cyan(fmt.Sprintf(`change path to: "%s"`, path))

			// fetch latest commit for all branch
			shell.GitFetch()
			// 按照分支创建 tag
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
					// 获取最新 COMMIT ID
					shell.GitCheckout().GetCommitId()
					ch <- fmt.Sprintf("%s - %s - commit id->%s", shell.Path, shell.CurrentBranch, shell.CommitId)
				}
			}
		}
		endCh <- true
	}(msgCh)

	for {
		select {
		case end, ok := <-endCh:
			if ok && end {
				colorlog.Success("Finished...")
				return
			}
		case msg, ok := <-msgCh:
			if ok {
				colorlog.Blue(msg)
				return
			}
		}
	}
}
