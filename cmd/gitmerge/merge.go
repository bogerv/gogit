package main

import (
	"fmt"
	"gitshell/config"
	"gitshell/internal/cmds"
	"gitshell/internal/constvar"
	"gitshell/pkg/colorlog"
	"log"
	"os"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const _splitFlag = "->"

func main() {
	// init and load config file
	config.Init("merge")
	conf := Config{}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("could not unmarshal config: %v", err)
	}

	//pathI := viper.Get("paths")
	paths := cast.ToStringSlice(conf.Paths) // 获取需要打 tag 的项目路径
	//flowsConfig := viper.Get("flows")
	flows := cast.ToStringSlice(conf.Flows) // 获取需要打 tag 的项目路径
	if len(conf.Flows) <= 0 {
		colorlog.Error(constvar.ErrNoFlow.Error())
		return
	}
	log.Printf("config: %v", conf)

	f, _ := os.OpenFile("./cmdmerge.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	for _, path := range paths {
		shell := new(cmds.Cmd)
		shell.Path = path
		// 切换执行 git 命令目录
		shell.Chdir()
		// change path and fetch
		colorlog.Cyan(fmt.Sprintf(`change path to: "%s"`, path))
		shell.GitFetch()

		for _, flow := range flows {
			if strings.Index(flow, _splitFlag) < 0 {
				continue
			}
			branches := strings.Split(flow, _splitFlag)
			if len(branches) != 2 {
				continue
			}

			branchForMerge := branches[0]
			checkoutTarget := branches[1]
			colorlog.Blue(fmt.Sprintf("git flow: %s-%s", branchForMerge, checkoutTarget))
			// 切换到要被合并的分支
			shell.CurrentBranch = branchForMerge
			shell.GitCheckout()

			// 拉取最新代码
			shell.GitPull()

			// 切换到分支
			shell.CurrentBranch = checkoutTarget
			shell.GitCheckout()

			// 拉取最新代码
			shell.GitPull()

			// 合并分支
			shell.GitMerge(branchForMerge)

			// 推送最新合并的内容
			shell.GitPushOrigin(checkoutTarget)

			// 获取最新 COMMIT ID
			shell.CurrentBranch = checkoutTarget
			shell.GetCommitId()
			colorlog.Success(fmt.Sprintf("%s--%s--commit id-> %s", shell.Path, checkoutTarget, shell.CommitId))

			log.SetOutput(f)
			if checkoutTarget != "test" {
				log.Printf("%s--%s-> %s", shell.Path, checkoutTarget, shell.CommitId)
			}
		}
		log.Printf("\n")
	}
	colorlog.Success("Merge Finished...")
}
