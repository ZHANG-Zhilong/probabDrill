package main

import (
	"fmt"
	"github.com/spf13/viper"
	"probabDrill/apps/probDrill/router"
	"probabDrill/cmd"
	probabDrill "probabDrill/conf"
)

var Conf = new(probabDrill.Config)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
	viper.SetConfigFile("./conf/config.yaml") // 指定配置文件
	err = viper.ReadInConfig()                // 读取配置信息
	if err != nil {                           // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err = viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}
	viper.WatchConfig()
	router.InitRouter()
}
