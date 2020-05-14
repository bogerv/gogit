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
	// 配置文件名称
	viper.SetConfigName("deltag")

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
	//branchesI := viper.Get("branches")
	//branches := cast.ToStringSlice(branchesI)
	//if len(branches) <= 0 {
	//	colorlog.Error(shellconst.ErrNoBranch.Error())
	//	return
	//}
	// 获取需要打 tag 的分支
	tagsI := viper.Get("tags")
	tags := cast.ToStringSlice(tagsI)
	if len(tags) <= 0 {
		colorlog.Error(shellconst.ErrNoBranch.Error())
		return
	}

	// tag 配置项
	//tagVersion := viper.GetString("tag")

	for _, path := range paths {
		shell := new(cmds.Cmd)
		shell.Path = path
		// 切换执行 git 命令目录
		shell.Chdir()
		// change path
		colorlog.Cyan(fmt.Sprintf(`change path to: "%s"`, path))

		// 按照分支创建 tag
		//for _, branch := range branches {
		//	// 按照分支创建 tag
		//	tagName := fmt.Sprintf("tag-%s-%s", branch, tagVersion)
		//	// remove origin tag
		//	colorlog.Yellow("remove tag: " + tagName)
		//	shell.PushOriginDelete(tagName)
		//}
		for _, tag := range tags {
			// remove origin tag
			colorlog.Yellow("remove tag: " + tag)
			shell.PushOriginDelete(tag)
		}
		fmt.Print("\n")
	}
	colorlog.Success("Finished...")
}
