package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"gitshell/cmds"
	"gitshell/colorlog"
)

const _splitFlag = "->"

func main() {
	// 配置文件名称
	viper.SetConfigName("merge")

	// 查找配置文件查的路径, 可以配置多个
	viper.AddConfigPath(".")
	viper.AddConfigPath("./resource")
	viper.AddConfigPath("../resource")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 获取需要打 tag 的项目路径
	flowsConfig := viper.Get("flows")
	flows := cast.ToStringSlice(flowsConfig)

	// 获取需要打 tag 的项目路径
	pathI := viper.Get("paths")
	paths := cast.ToStringSlice(pathI)

	f, _ := os.OpenFile("cmdmerge.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
