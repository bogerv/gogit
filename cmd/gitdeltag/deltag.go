package main

import (
	"fmt"
	"log"

	"gogit/config"
	"gogit/internal/cmds"
	"gogit/internal/constvar"
	"gogit/pkg/colorlog"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func main() {
	// init and load config file
	config.Init("deltag")
	conf := Config{}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("could not unmarshal config: %v", err)
	}

	// parse config
	conf.Paths = cast.ToStringSlice(conf.Paths) // projects path config
	conf.Tags = cast.ToStringSlice(conf.Tags)   // branches config
	if len(conf.Tags) <= 0 {
		colorlog.Error(constvar.ErrNoBranch.Error())
		return
	}
	log.Printf("config: %v", conf)

	for _, path := range conf.Paths {
		shell := new(cmds.Cmd)
		shell.Path = path
		// 切换执行 git 命令目录
		shell.Chdir()
		// change path
		colorlog.Cyan(fmt.Sprintf(`change path to: "%s"`, path))

		for _, tag := range conf.Tags {
			// remove origin tag
			colorlog.Yellow("remove tag: " + tag)
			shell.PushOriginDelete(tag)
		}
		fmt.Print("\n")
	}
	colorlog.Success("Finished...")
}
