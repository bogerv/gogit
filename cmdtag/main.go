package main

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"gitshell/cmds"
	"gitshell/colorlog"
	"gitshell/shellconst"
)

func main() {
	//defer shellutils.CatchPanic()
	//var pathCh = make(chan string, 1)
	//var endCh = make(chan bool)

	// 配置文件名称
	viper.SetConfigName("config")

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
	pathI := viper.Get("paths")
	paths := cast.ToStringSlice(pathI)

	// 获取需要打 tag 的分支
	branchesI := viper.Get("branches")
	branches := cast.ToStringSlice(branchesI)
	if len(branches) <= 0 {
		colorlog.Error(shellconst.ErrNoBranch.Error())
		return
	}

	// tag 配置项
	tagVersion := viper.GetString("tag")
	tagMessage := viper.GetString("message")

	for _, path := range paths {
		shell := new(cmds.Cmd)
		shell.Path = path
		// 切换执行 git 命令目录
		shell.Chdir()

		// change path and fetch
		colorlog.Cyan(fmt.Sprintf(`change path to: "%s"`, path))
		shell.GitFetch()
		// 按照分支创建 tag
		for _, branch := range branches {
			tagName := fmt.Sprintf("tag-%s-%s", branch, tagVersion)
			basedBranch := fmt.Sprintf("origin/%s", branch)

			// remove local same tag
			shell.DeleteTag(tagName)

			// add tag
			shell.AddTag(tagName, basedBranch, tagMessage)

			// push tag
			shell.GitPushOrigin(tagName)

			// 获取最新 COMMIT ID
			shell.GetCommitId(branch)
		}
	}

	//for {
	//	select {
	//	case end, ok := <-endCh:
	//		if ok && end {
	//			colorlog.Success("Finished...")
	//			return
	//		}
	//	}
	//}
}
