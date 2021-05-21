package main

import (
	"probabDrill/apps/probDrill/router"
	"probabDrill/cmd"
	"probabDrill/conf"
	_ "probabDrill/docs"
)

// @title probabDrill generation
// @version 1.0
// @license.name geit license.
// @host 171.16.1.106:4399
// @BasePath /v1

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
	conf.LoadConfig(cmd.Path)
	router.InitRouter()
}
