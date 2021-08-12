package main

import (
	"fmt"

	"coltcn.com/file-sync-server/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var conf config.Server

func main() {
	// 读取配置文件
	fmt.Println("读取配置文件")
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败！"))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件修改：", e.Name)
		if err := v.Unmarshal(&conf); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&conf); err != nil {
		fmt.Println(err)
	}
}
