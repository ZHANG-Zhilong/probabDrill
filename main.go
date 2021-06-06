package main

import (
	"probabDrill/apps/probDrill/router"
	"probabDrill/cmd"
	"probabDrill/conf"
	_ "probabDrill/docs"
)

// @title 基于沉积序列的三维地层概率模型研究
// @version 1.0
// @license.name  中山大学岩土工程与信息技术研究中心
// @host 171.16.1.107:4399
// @BasePath /v1

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
	conf.LoadConfig(cmd.Path)
	router.InitRouter()
}
