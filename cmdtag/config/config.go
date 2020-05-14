package config

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gitshell/colorlog"
	"gitshell/shellconst"
	"log"
)

type TagConfig struct {
	Paths    []string `json:"paths"`    // project absolute path (e.g. D:\golang\otc)
	Branches []string `json:"branches"` // project branches
	Version  string   `json:"version"`  // tag version
	Message  string   `json:"message"`  // tag message
}

func Init() *TagConfig {
	// set config file name
	viper.SetConfigName("config")

	// find config file from relative path
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("./resource")
	viper.AddConfigPath("../resource")
	viper.AddConfigPath("../../resource")

	// read context from config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// projects path config
	pathObj := viper.Get("paths")
	paths := cast.ToStringSlice(pathObj)

	// branches config
	branchesObj := viper.Get("branches")
	branches := cast.ToStringSlice(branchesObj)
	if len(branches) <= 0 {
		colorlog.Error(shellconst.ErrNoBranch.Error())
		log.Fatal()
	}

	// tag config
	version := viper.GetString("tag")
	message := viper.GetString("message")

	return &TagConfig{
		Paths:    paths,
		Branches: branches,
		Version:  version,
		Message:  message,
	}
}
