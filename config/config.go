package config

import (
	"github.com/spf13/viper"
	"log"
)

var (
	// config file name
	configName = "config"
	// config file paths
	configPaths = []string{
		"./",
		"./config",
		"../config",
	}
)

func Init(name string, path ...string) {
	if len(name) > 0 {
		configName = name
	}
	// 配置文件名称
	viper.SetConfigName(configName)
	// 查找配置文件查的路径, 可以配置多个
	configPaths = append(configPaths, path...)
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("load config err: %v", err)
	}

	//if err := viper.Unmarshal(&Conf); err != nil {
	//	log.Fatalf("could not unmarshal config: %v", err)
	//}
	//
	//log.Printf("config: %v", Conf)

	//viper.WatchConfig()
	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	if err := viper.Unmarshal(&Conf); err != nil {
	//		log.Printf("could not unmarshal config after changed: %v\n", err)
	//	}
	//	log.Printf("config: %v", Conf)
	//})
}
