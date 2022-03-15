package main

import (
	"fmt"
	"log"
	"strings"

	"gogit/config"
	"gogit/internal/cmds"
	"gogit/internal/constvar"
	"gogit/pkg/colorlog"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func main() {
	var msgCh = make(chan string, 1)
	var endCh = make(chan struct{})
	// init and load config file
	config.Init("addtag")
	conf := Config{}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("could not unmarshal config: %v", err)
	}

	// parse config
	conf.Paths = cast.ToStringSlice(conf.Paths)       // projects path config
	conf.Branches = cast.ToStringSlice(conf.Branches) // branches config
	if len(conf.Branches) <= 0 {
		colorlog.Error(constvar.ErrNoBranch.Error())
		log.Fatal()
	}
	log.Printf("config: %v", conf)

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
