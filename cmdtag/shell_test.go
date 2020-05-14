package main

import (
	"os"
	"os/exec"
	"testing"
)

// 查找配置文件查的路径, 可以配置多个
//viper.AddConfigPath("/etc/appname/")
//viper.AddConfigPath("$HOME/.appname")

func TestShell(t *testing.T) {
	cmd := exec.Command("git", "--help")

	//cmd := exec.Command("dir", ".")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
